package ocrcommon

import (
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const bridgeResponse = `{
			"telemetry":{
				"data_source":"data_source_test",
				"provider_requested_protocol":"provider_requested_protocol_test",
				"provider_requested_timestamp":922337203685477600,
				"provider_received_timestamp":-922337203685477600,
				"provider_data_stream_established":1,
				"provider_data_received":123456789,
				"provider_indicated_time":-123456789
			}
		}`

var trrs = pipeline.TaskRunResults{
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "ds1_parse", nil, nil, 1),
		},
		Result: pipeline.Result{
			Value: "123456",
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			BaseTask: pipeline.NewBaseTask(0, "ds2", nil, nil, 0),
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "ds2_parse", nil, nil, 1),
		},
		Result: pipeline.Result{
			Value: "12345678",
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			BaseTask: pipeline.NewBaseTask(0, "ds3", nil, nil, 0),
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "ds3_parse", nil, nil, 1),
		},
		Result: pipeline.Result{
			Value: "1234567890",
		},
	},
}

func TestShouldCollectTelemetry(t *testing.T) {
	j := job.Job{
		OCROracleSpec:  &job.OCROracleSpec{CaptureEATelemetry: true},
		OCR2OracleSpec: &job.OCR2OracleSpec{CaptureEATelemetry: true},
	}

	j.Type = job.Type(pipeline.OffchainReportingJobType)
	assert.True(t, shouldCollectTelemetry(&j))
	j.OCROracleSpec.CaptureEATelemetry = false
	assert.False(t, shouldCollectTelemetry(&j))

	j.Type = job.Type(pipeline.OffchainReporting2JobType)
	assert.True(t, shouldCollectTelemetry(&j))
	j.OCR2OracleSpec.CaptureEATelemetry = false
	assert.False(t, shouldCollectTelemetry(&j))

	j.Type = job.Type(pipeline.VRFJobType)
	assert.False(t, shouldCollectTelemetry(&j))
}

func TestGetContract(t *testing.T) {
	j := job.Job{
		OCROracleSpec:  &job.OCROracleSpec{CaptureEATelemetry: true},
		OCR2OracleSpec: &job.OCR2OracleSpec{CaptureEATelemetry: true},
	}
	contractAddress := ethkey.EIP55Address(utils.RandomAddress().String())

	j.Type = job.Type(pipeline.OffchainReportingJobType)
	j.OCROracleSpec.ContractAddress = contractAddress
	assert.Equal(t, contractAddress.String(), getContract(&j))

	j.Type = job.Type(pipeline.OffchainReporting2JobType)
	j.OCR2OracleSpec.ContractID = contractAddress.String()
	assert.Equal(t, contractAddress.String(), getContract(&j))

	j.Type = job.Type(pipeline.VRFJobType)
	assert.Empty(t, getContract(&j))
}

func TestGetChainID(t *testing.T) {
	j := job.Job{
		OCROracleSpec:  &job.OCROracleSpec{CaptureEATelemetry: true},
		OCR2OracleSpec: &job.OCR2OracleSpec{CaptureEATelemetry: true},
	}

	j.Type = job.Type(pipeline.OffchainReportingJobType)
	j.OCROracleSpec.EVMChainID = (*utils.Big)(big.NewInt(1234567890))
	assert.Equal(t, "1234567890", getChainID(&j))

	j.Type = job.Type(pipeline.OffchainReporting2JobType)
	j.OCR2OracleSpec.RelayConfig = make(map[string]interface{})
	j.OCR2OracleSpec.RelayConfig["chainID"] = "foo"
	assert.Equal(t, "foo", getChainID(&j))

	j.Type = job.Type(pipeline.VRFJobType)
	assert.Empty(t, getChainID(&j))
}

func TestParseEATelemetry(t *testing.T) {
	ea, err := parseEATelemetry([]byte(bridgeResponse))
	assert.NoError(t, err)
	assert.Equal(t, ea.DataSource, "data_source_test")
	assert.Equal(t, ea.ProviderRequestedProtocol, "provider_requested_protocol_test")
	assert.Equal(t, ea.ProviderRequestedTimestamp, int64(922337203685477600))
	assert.Equal(t, ea.ProviderReceivedTimestamp, int64(-922337203685477600))
	assert.Equal(t, ea.ProviderDataStreamEstablished, int64(1))
	assert.Equal(t, ea.ProviderDataReceived, int64(123456789))
	assert.Equal(t, ea.ProviderIndicatedTime, int64(-123456789))

	_, err = parseEATelemetry(nil)
	assert.Error(t, err)
}

func TestGetJsonParsedValue(t *testing.T) {

	resp := getJsonParsedValue(trrs[0], &trrs)
	assert.Equal(t, "123456", resp.String())

	trrs[1].Result.Value = nil
	resp = getJsonParsedValue(trrs[0], &trrs)
	assert.Nil(t, resp)

	resp = getJsonParsedValue(trrs[1], &trrs)
	assert.Nil(t, resp)

}

func TestSendEATelemetry(t *testing.T) {
	wg := sync.WaitGroup{}
	ingressClient := mocks.NewTelemetryIngressClient(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("0xa", synchronization.EnhancedEA)

	var sentMessage []byte
	ingressClient.On("Send", mock.AnythingOfType("synchronization.TelemPayload")).Return().Run(func(args mock.Arguments) {
		sentMessage = args[0].(synchronization.TelemPayload).Telemetry
		wg.Done()
	})

	feedAddress := utils.RandomAddress()

	ds := inMemoryDataSource{
		jb: job.Job{
			Type: job.Type(pipeline.OffchainReportingJobType),
			OCROracleSpec: &job.OCROracleSpec{
				ContractAddress:    ethkey.EIP55AddressFromAddress(feedAddress),
				CaptureEATelemetry: true,
				EVMChainID:         (*utils.Big)(big.NewInt(9)),
			},
		},
		lggr:               nil,
		monitoringEndpoint: monitoringEndpoint,
	}
	trrs := pipeline.TaskRunResults{
		pipeline.TaskRunResult{
			Task: &pipeline.BridgeTask{
				BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			},
			Result: pipeline.Result{
				Value: bridgeResponse,
			},
		},
		pipeline.TaskRunResult{
			Task: &pipeline.JSONParseTask{
				BaseTask: pipeline.NewBaseTask(1, "ds1", nil, nil, 1),
			},
			Result: pipeline.Result{
				Value: "1234567890",
			},
		},
	}
	fr := pipeline.FinalResult{
		Values:      []interface{}{"123456"},
		AllErrors:   nil,
		FatalErrors: []error{nil},
	}

	wg.Add(1)
	collectEATelemetry(&ds, &trrs, &fr)

	expectedTelemetry := telem.EnhancedEA{
		DataSource:                    "data_source_test",
		Value:                         1234567890,
		BridgeTaskRunStartedTimestamp: trrs[0].CreatedAt.UnixMilli(),
		BridgeTaskRunEndedTimestamp:   trrs[0].FinishedAt.Time.UnixMilli(),
		ProviderRequestedProtocol:     "provider_requested_protocol_test",
		ProviderRequestedTimestamp:    922337203685477600,
		ProviderReceivedTimestamp:     -922337203685477600,
		ProviderDataStreamEstablished: 1,
		ProviderDataReceived:          123456789,
		ProviderIndicatedTime:         -123456789,
		Feed:                          feedAddress.String(),
		ChainId:                       "9",
		Observation:                   123456,
		Round:                         0,
		Epoch:                         0,
	}

	expectedMessage, _ := proto.Marshal(&expectedTelemetry)
	wg.Wait()
	assert.Equal(t, expectedMessage, sentMessage)
}

func TestGetObservation(t *testing.T) {
	lggr, logs := logger.TestLoggerObserved(t, zap.WarnLevel)
	ds := &inMemoryDataSource{
		jb: job.Job{
			ID:   1234567890,
			Type: job.Type(pipeline.OffchainReportingJobType),
			OCROracleSpec: &job.OCROracleSpec{
				CaptureEATelemetry: true,
			},
		},
		lggr: lggr,
	}

	obs := getObservation(ds, &pipeline.FinalResult{})
	assert.Equal(t, obs, int64(0))
	assert.Equal(t, logs.Len(), 1)
	assert.Contains(t, logs.All()[0].Message, "cannot get singular result")

	finalResult := &pipeline.FinalResult{
		Values:      []interface{}{"123456"},
		AllErrors:   nil,
		FatalErrors: []error{nil},
	}
	obs = getObservation(ds, finalResult)
	assert.Equal(t, obs, int64(123456))
}

func TestCollectAndSend(t *testing.T) {
	ingressClient := mocks.NewTelemetryIngressClient(t)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("0xa", synchronization.EnhancedEA)
	ingressClient.On("Send", mock.AnythingOfType("synchronization.TelemPayload")).Return()

	lggr, logs := logger.TestLoggerObserved(t, zap.WarnLevel)
	ds := &inMemoryDataSource{
		jb: job.Job{
			ID:   1234567890,
			Type: job.Type(pipeline.OffchainReportingJobType),
			OCROracleSpec: &job.OCROracleSpec{
				CaptureEATelemetry: true,
			},
		},
		lggr:               lggr,
		monitoringEndpoint: monitoringEndpoint,
	}

	finalResult := &pipeline.FinalResult{
		Values:      []interface{}{"123456"},
		AllErrors:   nil,
		FatalErrors: []error{nil},
	}

	badTrrs := &pipeline.TaskRunResults{
		pipeline.TaskRunResult{
			Task: &pipeline.BridgeTask{
				BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			},
			Result: pipeline.Result{
				Value: nil,
			},
		}}

	collectAndSend(ds, badTrrs, finalResult)
	assert.Contains(t, logs.All()[0].Message, "cannot get bridge response from bridge task")

	badTrrs = &pipeline.TaskRunResults{
		pipeline.TaskRunResult{
			Task: &pipeline.BridgeTask{
				BaseTask: pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			},
			Result: pipeline.Result{
				Value: "[]",
			},
		}}
	collectAndSend(ds, badTrrs, finalResult)
	assert.Equal(t, logs.Len(), 3)
	assert.Contains(t, logs.All()[1].Message, "cannot parse EA telemetry")
	assert.Contains(t, logs.All()[2].Message, "cannot get json parse value")
}

func BenchmarkCollectEATelemetry(b *testing.B) {
	wg := sync.WaitGroup{}
	ingressClient := mocks.NewTelemetryIngressClient(b)
	ingressAgent := telemetry.NewIngressAgentWrapper(ingressClient)
	monitoringEndpoint := ingressAgent.GenMonitoringEndpoint("0xa", synchronization.EnhancedEA)

	ingressClient.On("Send", mock.AnythingOfType("synchronization.TelemPayload")).Return().Run(func(args mock.Arguments) {
		wg.Done()
	})

	ds := inMemoryDataSource{
		jb: job.Job{
			Type: job.Type(pipeline.OffchainReportingJobType),
			OCROracleSpec: &job.OCROracleSpec{
				ContractAddress:    ethkey.EIP55AddressFromAddress(utils.RandomAddress()),
				CaptureEATelemetry: true,
				EVMChainID:         (*utils.Big)(big.NewInt(9)),
			},
		},
		lggr:               nil,
		monitoringEndpoint: monitoringEndpoint,
	}
	finalResult := pipeline.FinalResult{
		Values:      []interface{}{"123456"},
		AllErrors:   nil,
		FatalErrors: []error{nil},
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		wg.Add(1)
		collectEATelemetry(&ds, &trrs, &finalResult)
	}
	wg.Wait()
}
