package query

import (
	"fas/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertFlowLog constructs a insertion template for a single flow-log
func InsertFlowLog(log data.FlowLog) bson.D {
	return bson.D{
		{"src_app", log.Src_App},
		{"dest_app", log.Dest_App},
		{"vpc_id", log.Vpc_Id},
		{"bytes_tx", log.Bytes_Tx},
		{"bytes_rx", log.Bytes_Rx},
		{"hour", log.Hour},
	}
}

// AggregationByHourPipeline constructs a aggregation query pipeline to get network stats data by hour
func AggregationByHourPipeline(hour int) mongo.Pipeline {
	// provision each aggregation stage
	matchStage := bson.D{{"$match", bson.D{{"hour", hour}}}}
	groupStage := bson.D{
		{"$group",
			bson.D{
				{"_id",
					bson.D{
						{"src_app", "$src_app"},
						{"dest_app", "$dest_app"},
						{"vpc_id", "$vpc_id"},
						{"hour", "$hour"},
					},
				},
				{"bytes_tx",
					bson.D{
						{"$sum", "$bytes_tx"},
					},
				},
				{"bytes_rx",
					bson.D{
						{"$sum", "$bytes_rx"},
					},
				},
			},
		},
	}
	projectStage := bson.D{
		{
			"$project",
			bson.D{
				{"_id", 0},
				{"src_app", "$_id.src_app"},
				{"dest_app", "$_id.dest_app"},
				{"vpc_id", "$_id.vpc_id"},
				{"hour", "$_id.hour"},
				{"bytes_tx", "$bytes_tx"},
				{"bytes_rx", "$bytes_rx"},
			},
		},
	}

	// assemble stages to the aggregation pipeline
	return mongo.Pipeline{matchStage, groupStage, projectStage}
}
