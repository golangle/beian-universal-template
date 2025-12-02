package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplitLineInfo(t *testing.T) {
	line := "标题|host1.com,host2.com|版权信息|备案号"
	lineSegment := strings.Split(line, "|") //分割标题、主机名、版权信息和备案号
	for idx, val := range lineSegment {
		fmt.Printf("idx=%d,val=%s\n", idx, val)
	}
}
