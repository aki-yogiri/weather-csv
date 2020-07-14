package handler

import (
	"context"
	"errors"
	pb "github.com/aki-yogiri/weather-store/pb/weather"
	"github.com/golang/protobuf/ptypes"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"time"
)

type StoreServerEnv struct {
	Host string
	Port string
}

func DownloadWeatherCSV(env StoreServerEnv) echo.HandlerFunc {
	return func(c echo.Context) error {
		location := c.QueryParam("location")
		if location == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "need query param: location")
		}

		query := &pb.QueryMessage{}
		query.Location = location

		layout := "2006-01-02T15:04:05Z"
		dtstart := c.QueryParam("dtstart")
		if dtstart != "" {
			t, err := time.Parse(layout, dtstart)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid format dtstart: "+dtstart)
			}

			query.DatetimeStart, err = ptypes.TimestampProto(t)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid format dtstart: "+dtstart)
			}
		}
		dtend := c.QueryParam("dtend")
		if dtend != "" {
			t, err := time.Parse(layout, dtend)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid format dtend: "+dtend)
			}

			query.DatetimeEnd, err = ptypes.TimestampProto(t)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid format dtend: "+dtend)
			}
		}

		resp, err := executeQuery(env, query)
		if err != nil {
			log.Printf("Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		data, err := makeCSV(resp)

		if err != nil {
			log.Printf("Error %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename=weather-"`+time.Now().Format(layout)+`.csv"`)

		return c.Blob(http.StatusOK, "text/csv", data)

	}
}

func executeQuery(env StoreServerEnv, query *pb.QueryMessage) (*pb.WeatherReply, error) {
	conn, err := grpc.Dial(env.Host+":"+env.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
		return nil, errors.New("can not connection store")
	}
	defer conn.Close()

	client := pb.NewWeatherClient(conn)
	res, err := client.GetWeather(context.TODO(), query)
	if err != nil {
		log.Printf("Error %v", err)
		return nil, err
	}

	return res, nil
}

func makeCSV(weather *pb.WeatherReply) ([]byte, error) {

	content := make([]byte, 0)
	encode := "\n"
	separate := ","

	content = append(content, "timestamp,location,weather,temperature,clouds,wind,wind_deg"...)
	content = append(content, encode...)
	for _, v := range weather.Weather {
		content = append(content, ptypes.TimestampString(v.Timestamp)...)
		content = append(content, separate...)
		content = append(content, v.Location...)
		content = append(content, separate...)
		content = append(content, v.Weather...)
		content = append(content, separate...)
		content = append(content, strconv.FormatFloat(v.Temperature, 'f', 4, 64)...)
		content = append(content, separate...)
		content = append(content, strconv.FormatUint(uint64(v.Clouds), 10)...)
		content = append(content, separate...)
		content = append(content, strconv.FormatFloat(v.Wind, 'f', 4, 64)...)
		content = append(content, separate...)
		content = append(content, strconv.FormatUint(uint64(v.WindDeg), 10)...)

		content = append(content, encode...)
	}

	return content, nil

}
