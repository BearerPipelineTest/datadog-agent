// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !kubeapiserver
// +build !kubeapiserver

package hostinfo

import "context"

func apiserverNodeLabels(ctx context.Context, nodeName string) (map[string]string, error) {
	return nil, nil
}

func apiserverNodeAnnotations(ctx context.Context, nodeName string) (map[string]string, error) {
	return nil, nil
}
