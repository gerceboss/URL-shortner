package routes

import (
	"github.com/gerceboss/url-shortner/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)
func ResolveURL(c *fiber.Ctx) error{
	url:=c.Params("url")
	r:=database.CreateClient(0)
	defer r.Close()
	//key value pair database
	value,err:= r.Get(database.Ctx,url).Result()

	if err==redis.Nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":"not found in the database"})
	}else if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"cannot connect to the database"})
	}

	redisInc:=database.CreateClient(1)
	defer redisInc.Close()

	_ =redisInc.Incr(database.Ctx,"counter")
	return c.Redirect(value,301)
}