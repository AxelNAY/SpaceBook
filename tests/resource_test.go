package tests

import (
	"net/http"
	"spacebook/config"
	"spacebook/models"
	"spacebook/handlers"
    "log"
    "os"
    "testing"

	"github.com/labstack/echo/v4"
)

func TestCreateResource(t *testing.T) {
	newRoom := model.resource{
                Username: "yemiwebby",
                Password: "test",
        }
}