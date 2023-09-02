package dynamo_container

import (
	"context"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/walteh/testrc/pkg/dynamo"
	"github.com/walteh/testrc/tests/containers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog"
)

func NewMockAPI(t *testing.T, ctx context.Context, cfg aws.Config, input *dynamodb.CreateTableInput) dynamo.DynamoDBAPIProvisioner {
	return containers.WrapTestClient(t, ctx, newMockClient(cfg, containers.GetHttp(global)), input)
}

func PrintScanAsTable(ctx context.Context, svc dynamo.DynamoDBAPIProvisioner, tbl string) {
	zerolog.Ctx(ctx).Info().Msg("Scanning table " + tbl)

	input := &dynamodb.ScanInput{TableName: aws.String(tbl)}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := map[string]int{}
	rows := make([]map[string]interface{}, 0)

	paginator := dynamodb.NewScanPaginator(svc, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			fmt.Printf("failed to get page: %v", err)
			return
		}

		if len(page.Items) == 0 {
			continue
		}

		for _, item := range page.Items {
			row := map[string]interface{}{}
			for i, h := range item {
				header[i]++
				var marsh interface{}
				err := attributevalue.Unmarshal(h, &marsh)
				if err != nil {
					row[i] = h
				}
				switch m := marsh.(type) {
				case []byte:
					row[i] = string(m)
				case []interface{}:
					rat := make([]string, 0)
					for _, v := range m {
						switch vv := v.(type) {
						case float64:
							rat = append(rat, fmt.Sprintf("%f", vv))
						default:
							rat = append(rat, fmt.Sprintf("%v", vv))
						}
					}
					row[i] = rat
				case float64:
					row[i] = fmt.Sprintf("%f", m)
				default:
					row[i] = marsh
				}
			}
			rows = append(rows, row)
		}
	}

	str := make([]string, 0)
	for k := range header {
		str = append(str, k)
	}

	sort.Strings(str)

	tmp := make([]interface{}, len(str))
	for i, v := range str {
		tmp[i] = v
	}

	t.AppendHeader(table.Row(tmp))

	for _, d := range rows {
		row := make([]interface{}, len(str))
		for i, v := range str {
			row[i] = d[v]
		}
		t.AppendRow(row)
	}

	t.SetTitle(fmt.Sprintf("Table Data: %s", tbl))

	t.Render()
}

func PrintTableCounts(ctx context.Context, svc dynamo.DynamoDBAPIProvisioner, tbl string) {
	zerolog.Ctx(ctx).Info().Msg("Scanning table " + tbl)

	input := &dynamodb.ScanInput{
		TableName: aws.String(tbl),
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := make(map[string]int, 0)
	headerNull := make(map[string]int, 0)

	paginator := dynamodb.NewScanPaginator(svc, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			fmt.Printf("failed to get page: %v", err)
			return
		}

		if len(page.Items) == 0 {
			continue
		}

		for _, item := range page.Items {
			for i, h := range item {
				var marsh interface{}
				err := attributevalue.Unmarshal(h, &marsh)
				if err != nil {
					headerNull[i]++
				}
				header[i]++
			}

		}
	}

	t.SetTitle(fmt.Sprintf("Table Counts: %s", tbl))

	t.AppendHeader(table.Row{"Attribute", "Count", "Null Count"})

	str := make([]string, 0)
	for k := range header {
		str = append(str, k)
	}
	sort.Strings(str)
	for _, k := range str {
		t.AppendRow([]interface{}{k, header[k], headerNull[k]})
	}

	t.Render()
}
