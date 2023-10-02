package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gerceboss/url-shortner/database"
	"github.com/gerceboss/url-shortner/helper"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type request struct{
	URL	 string `json:"url"`
	CustomShort string `json:"short"`
	ExpiryTime time.Duration `json:"expiry"`
}
type response struct{
	URL string `json:"url"`
	CustomShort string `json:"short"`
	ExpiryTime time.Duration `json:"expiry"`
	RateLimitReset time.Duration `json:"rate_limit_reset"`
	RateRemaining int `json:"rate_limit"`
}
func ShortenURL(c  *fiber.Ctx) error{
	body:=new(request)
	if err:=c.BodyParser(&body);err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"can't parse json",})
	}

	//check rateLimit using the user's stats in db
	//check IP address if already in db
	r:=database.CreateClient(1)
	defer r.Close() //called at end of this function's call stack to disconnect the database

	v,err:=r.Get(database.Ctx, c.IP()).Result()
	if err== redis.Nil{
		//user not found in db
		// IP,quoateLeft,timeitResets
		_=r.Set(database.Ctx,c.IP(),os.Getenv("API_QUOTA"),30*60*time.Second).Err()	
	}else{
		
		vInt,_:=strconv.Atoi(v)
		if vInt<=0{
			limit,_:=r.TTL(database.Ctx,c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":"try after some time ",
				"rate_limit_reset":limit / time.Nanosecond / time.Minute,
			})
		}
	}
	//check if input is aactually a URL
	if  !govalidator.IsURL(body.URL){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid URL",})
	}

	//check for domain error so that there is not any infinite loop
	if !helper.RemoveDomainError(body.URL){
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error":"invlid domain"})
	}

	//enforce that url must be https or ssl if its not convert it 
	body.URL= helper.EnforceHTTP(body.URL)

	var id string
	if body.CustomShort==""{
		id=uuid.New().String()[:8]
	}else{
		id=body.CustomShort
	}

	r1:=database.CreateClient(1)
	defer r1.Close()
	//check if ID already in USE
	value,_:=r1.Get(database.Ctx,id).Result()
	if value!=""{
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error":"URL already in use"})
	}

	if body.ExpiryTime==0{
		body.ExpiryTime=1 
	}
	err=r1.Set(database.Ctx,id,body.URL,body.ExpiryTime*24*3600*time.Second).Err()

	if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"dunno",
		})
	}
	

	res:=response{
		URL :body.URL,
		CustomShort:os.Getenv("DOMAIN")+ "/" + id,
		ExpiryTime:body.ExpiryTime,
		RateLimitReset:30,
		RateRemaining:10,
	}
	r.Decr(database.Ctx,c.IP())
	value,_=r.Get(database.Ctx,c.IP()).Result()
	res.RateRemaining,_=strconv.Atoi(value)
	ttl,_:=r.TTL(database.Ctx,c.IP()).Result()
	res.RateLimitReset= ttl/time.Nanosecond /time.Minute

	return c.Status(fiber.StatusOK).JSON(res)
}