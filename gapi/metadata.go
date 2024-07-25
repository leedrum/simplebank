package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader        = "x-forwarded-for"
	userAgentHeader            = "user-agent"
)

type MetaData struct {
	ClientIP  string
	UserAgent string
}

func (server *Server) extractMetadata(ctx context.Context) *MetaData {
	metaData := &MetaData{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("metadata: %v", md)

		if userAgent := md.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
			metaData.UserAgent = userAgent[0]
		}

		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			metaData.UserAgent = userAgent[0]
		}

		if clientIP := md.Get(xForwardedForHeader); len(clientIP) > 0 {
			metaData.ClientIP = clientIP[0]
		}
	}

	if peerInfo, ok := peer.FromContext(ctx); ok {
		if peerInfo.Addr != nil {
			metaData.ClientIP = peerInfo.Addr.String()
		}
	}

	return metaData
}
