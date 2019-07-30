package ali_mns

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	decoder := NewAliMNSDecoder()

	mMNSClient := &mockMNSClient{}
	fresp := &fasthttp.Response{}
	fresp.Header.Add(headerKeyRequestID, "test-request-id")
	fresp.SetStatusCode(200)
	fresp.SetBody([]byte(`<Message xmlns="http://mns.aliyuncs.com/doc/v1"><MessageId>86E3874C5D49</MessageId><MessageBodyMD5>0F0479874BF6F4A7281099B1</MessageBodyMD5></Message>`))
	mMNSClient.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fresp, nil)

	resp := &MessageSendResponse{}
	sc, err := send(mMNSClient, decoder, POST, map[string]string{}, MessageSendRequest{}, "queues", resp)
	assert.Nil(t, err)
	assert.Equal(t, 200, sc)
	assert.Equal(t, "test-request-id", resp.RequestID)
}


func TestSendError(t *testing.T) {
	decoder := NewAliMNSDecoder()

	mMNSClient := &mockMNSClient{}
	fresp := &fasthttp.Response{}
	fresp.Header.Add(headerKeyRequestID, "test-request-id")
	fresp.SetStatusCode(400)
	fresp.SetBody([]byte(`<Error xmlns="http://mns.aliyuncs.com/doc/v1"></Error>`))
	mMNSClient.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fresp, nil)

	resp := &MessageSendResponse{}
	sc, err := send(mMNSClient, decoder, POST, map[string]string{}, MessageSendRequest{}, "queues", resp)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request id: test-request-id")
	assert.Equal(t, 400, sc)
	assert.Equal(t, "", resp.RequestID)
}
