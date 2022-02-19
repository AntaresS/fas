package server

import (
	"context"
	"encoding/json"
	"fas/data"
	"fas/data/query"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

// ReportNetworkFlow is the POST endpoint to insert flow-log data
func (s *Server) ReportNetworkFlow(c echo.Context) error {
	c.Logger().SetLevel(log.DEBUG)
	msg := []data.FlowLog{}
	if err := c.Bind(&msg); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(msg); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	s.insertDocuments(msg)

	return nil
}

func (s *Server) insertDocuments(logs []data.FlowLog) {

	flowLogCollection := s.mongoClient.Database("monitoring").Collection("flowlog")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	for _, l := range logs {
		res, err := flowLogCollection.InsertOne(ctx, query.InsertFlowLog(l))
		if err != nil {
			s.logger.Errorf("failed to insert doc: %v with err: %v", l, err)
		}
		s.logger.Infof("inserted with id: %v", res.InsertedID)
	}
}

// AggregateFlowDataByHour is the GET endpoint to aggregate network stats by hour
func (s *Server) AggregateFlowDataByHour(c echo.Context) error {

	hour := c.QueryParam("hour")
	// query param validations
	if hour == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required parameter \"hour\" or the value is empty")
	}

	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad parameter value: %v, requires integer", hour))
	}

	flowLogCollection := s.mongoClient.Database("monitoring").Collection("flowlog")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// execute the aggregation based on the given hour
	cursor, err := flowLogCollection.Aggregate(ctx, query.AggregationByHourPipeline(hourInt))
	if err != nil {
		c.Logger().Errorf("query request failed with err: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// decode all results
	results := []data.FlowLog{}
	if err := cursor.All(ctx, &results); err != nil {
		c.Logger().Errorf("failed to fetch query results with err: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(results)
}
