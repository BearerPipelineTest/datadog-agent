// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf
// +build linux_bpf

package compiler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/DataDog/datadog-agent/pkg/util/kernel"
)

var (
	datadogAgentEmbeddedPath = "/opt/datadog-agent/embedded"
	clangBinPath             = filepath.Join(datadogAgentEmbeddedPath, "bin/clang")
	llcBinPath               = filepath.Join(datadogAgentEmbeddedPath, "bin/llc")
	embeddedIncludePath      = filepath.Join(datadogAgentEmbeddedPath, "include")
	defaultFlags             = []string{
		"-D__KERNEL__",
		"-DCONFIG_64BIT",
		"-D__BPF_TRACING__",
		`-DKBUILD_MODNAME="ddsysprobe"`,
		"-Wno-unused-value",
		"-Wno-pointer-sign",
		"-Wno-compare-distinct-pointer-types",
		"-Wunused",
		"-Wall",
		"-Werror",
		"-emit-llvm",
		"-O2",
		"-fno-stack-protector",
		"-fno-color-diagnostics",
		"-fno-unwind-tables",
		"-fno-asynchronous-unwind-tables",
		"-fno-jump-tables",
		"-nostdinc",
	}
)

func CompileToObjectFile(in io.Reader, outputFile string, cflags []string, headerDirs []string) error {
	if len(headerDirs) == 0 {
		fmt.Errorf("unable to find kernel headers")
	}

	arch := kernel.Arch()
	if arch == "" {
		fmt.Errorf("unable to get kernel arch for %s", runtime.GOARCH)
	}

	cflags = append(cflags, defaultFlags...)
	for _, d := range headerDirs {
		cflags = append(cflags,
			fmt.Sprintf("-isystem%s/arch/%s/include", d, arch),
			fmt.Sprintf("-isystem%s/arch/%s/include/generated", d, arch),
			fmt.Sprintf("-isystem%s/include", d),
			fmt.Sprintf("-isystem%s/arch/%s/include/uapi", d, arch),
			fmt.Sprintf("-isystem%s/arch/%s/include/generated/uapi", d, arch),
			fmt.Sprintf("-isystem%s/include/uapi", d),
			fmt.Sprintf("-isystem%s/include/generated/uapi", d),
		)
	}
	cflags = append(cflags, fmt.Sprintf("-isystem%s", embeddedIncludePath))
	cflags = append(cflags, "-c", "-x", "c", "-o", "-", "-")

	var clangOut, clangErr, llcOut, llcErr bytes.Buffer

	clangCtx, clangCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer clangCancel()

	compileToBC := exec.CommandContext(clangCtx, clangBinPath, cflags...)
	compileToBC.Stdin = in
	compileToBC.Stdout = &clangOut
	compileToBC.Stderr = &clangErr

	err := compileToBC.Run()

	if err != nil {
		if clangCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("error compiling asset to bytecode: operation timed out")
		}

		return fmt.Errorf("error compiling asset to bytecode: %s", clangErr.String())
	}

	llcCtx, llcCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer llcCancel()

	bcToObj := exec.CommandContext(llcCtx, llcBinPath, "-march=bpf", "-filetype=obj", "-o", outputFile, "-")
	bcToObj.Stdin = &clangOut
	bcToObj.Stdout = &llcOut
	bcToObj.Stderr = &llcErr

	err = bcToObj.Run()
	if err != nil {
		if llcCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("error compiling bytecode to object file: operation timed out")
		}

		return fmt.Errorf("error compiling bytecode to object file: %s", llcErr.String())
	}
	return nil
}
