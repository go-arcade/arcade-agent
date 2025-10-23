package grpc

import (
	"fmt"

	agentv1 "github.com/go-arcade/arcade-agent/api/agent/v1"
	pipelinev1 "github.com/go-arcade/arcade-agent/api/pipeline/v1"
	taskv1 "github.com/go-arcade/arcade-agent/api/task/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GrpcClient gRPC客户端封装
type GrpcClient struct {
	conn           *grpc.ClientConn
	agentClient    agentv1.AgentServiceClient
	taskClient     taskv1.TaskServiceClient
	pipelineClient pipelinev1.PipelineServiceClient
	serverAddr     string
}

// NewGrpcClient 创建新的gRPC客户端
func NewGrpcClient(serverAddr string) (*GrpcClient, error) {
	// 设置连接选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// 连接到gRPC服务器
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("连接gRPC服务器失败: %v", err)
	}

	return &GrpcClient{
		conn:           conn,
		agentClient:    agentv1.NewAgentServiceClient(conn),
		taskClient:     taskv1.NewTaskServiceClient(conn),
		pipelineClient: pipelinev1.NewPipelineServiceClient(conn),
		serverAddr:     serverAddr,
	}, nil
}

// Close 关闭gRPC连接
func (c *GrpcClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetAgentClient 获取Agent服务客户端
func (c *GrpcClient) GetAgentClient() agentv1.AgentServiceClient {
	return c.agentClient
}

// GetTaskClient 获取Task服务客户端
func (c *GrpcClient) GetTaskClient() taskv1.TaskServiceClient {
	return c.taskClient
}

// GetPipelineClient 获取Pipeline服务客户端
func (c *GrpcClient) GetPipelineClient() pipelinev1.PipelineServiceClient {
	return c.pipelineClient
}

// GetServerAddr 获取服务器地址
func (c *GrpcClient) GetServerAddr() string {
	return c.serverAddr
}
