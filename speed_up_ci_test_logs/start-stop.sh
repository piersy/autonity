#!/bin/bash

./speed_up_ci -test.run TestTendermintStopUpToFNodes -test.v -test.count 100 -test.timeout 120m  2>&1 | tee investigate_TestTendermintStopUpToFNodes.log

./speed_up_ci -test.run TestTendermintStartStopSingleNode -test.v -test.count 100 -test.timeout 120m  2>&1 | tee investigate_TestTendermintStartStopSingleNode.log

./speed_up_ci -test.run TestTendermintStartStopFNodes -test.v -test.count 100 -test.timeout 120m  2>&1 | tee investigate_TestTendermintStartStopFNodes.log

./speed_up_ci -test.run TestTendermintStartStopFPlusOneNodes -test.v -test.count 100 -test.timeout 120m  2>&1 | tee investigate_TestTendermintStartStopFPlusOneNodes.log

./speed_up_ci -test.run TestTendermintStartStopFPlusTwoNodes -test.v -test.count 100 -test.timeout 120m  2>&1 | tee investigate_TestTendermintStartStopFPlusTwoNodes.log

./speed_up_ci -test.run TestTendermintStartStopAllNodes -test.v -test.count 100 -test.timeout 120m  2>&1 | tee investigate_TestTendermintStartStopAllNodes.log