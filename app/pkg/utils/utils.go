package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thedatashed/xlsxreader"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var TraceIDHeader string = `x-trace-id`

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ParseMultipartFileXlsx(sheetName string, multipartFile *multipart.FileHeader) (result []string, err error) {
	file, err := multipartFile.Open()
	if err != nil {
		return []string{}, err
	}
	defer file.Close()
	var buf bytes.Buffer
	io.Copy(&buf, file)

	xl, _ := xlsxreader.NewReader(buf.Bytes())

	for row := range xl.ReadRows(sheetName) {
		if len(row.Cells) == 0 {
			continue
		}
		result = append(result, row.Cells[0].Value)
	}
	if len(result) == 0 {
		err = errors.New("empty file")
	}
	return
}

func StartFiberTrace(c *fiber.Ctx, spanName string) (context.Context, trace.Span) {
	ctx, span := otel.Tracer("").Start(c.Context(), spanName)
	c.Set(TraceIDHeader, span.SpanContext().TraceID().String())
	if clientID, ok := c.Locals("clientID").(int); ok {
		span.SetAttributes(attribute.Int("clientID", clientID))
	}
	if clientID, ok := c.Locals("APIID").(int); ok {
		span.SetAttributes(attribute.Int("APIID", clientID))
	}
	if userID, ok := c.Locals("userID").(int); ok {
		span.SetAttributes(attribute.Int("userID", userID))
	}
	return ctx, span
}

func BuildCookie(name string, value string, expires time.Time) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.HTTPOnly = true
	cookie.Expires = expires
	cookie.Path = "/"

	return cookie
}

func StructToPrettyJsonString(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}
