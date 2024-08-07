/*
Copyright 2020 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package observers

import (
	"fmt"
	"time"

	flaggerv1 "github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/fluxcd/flagger/pkg/metrics/providers"
)

var nginxQueries = map[string]string{
	"request-success-rate": `
	sum(
		rate(
			nginx_ingress_controller_requests{
				namespace="{{ namespace }}",
				ingress="{{ ingress }}",
				canary!="",
				status!~"5.*"
			}[{{ interval }}]
		)
	) 
	/ 
	sum(
		rate(
			nginx_ingress_controller_requests{
				namespace="{{ namespace }}",
				ingress="{{ ingress }}",
				canary!=""
			}[{{ interval }}]
		)
	) 
	* 100`,
	"request-duration": `
	sum(
		rate(
			nginx_ingress_controller_response_duration_seconds_sum{
				namespace="{{ namespace }}",
				ingress="{{ ingress }}",
				canary!=""
			}[{{ interval }}]
		)
	) 
	/ 
	sum(
		rate(
			nginx_ingress_controller_response_duration_seconds_count{
				namespace="{{ namespace }}",
				ingress="{{ ingress }}",
				canary!=""
			}[{{ interval }}]
		)
	) 
	* 1000`,
}

type NginxObserver struct {
	client providers.Interface
}

func (ob *NginxObserver) GetRequestSuccessRate(model flaggerv1.MetricTemplateModel) (float64, error) {
	query, err := RenderQuery(nginxQueries["request-success-rate"], model)
	if err != nil {
		return 0, fmt.Errorf("rendering query failed: %w", err)
	}

	value, err := ob.client.RunQuery(query)
	if err != nil {
		return 0, fmt.Errorf("running query failed: %w", err)
	}

	return value, nil
}

func (ob *NginxObserver) GetRequestDuration(model flaggerv1.MetricTemplateModel) (time.Duration, error) {
	query, err := RenderQuery(nginxQueries["request-duration"], model)
	if err != nil {
		return 0, fmt.Errorf("rendering query failed: %w", err)
	}

	value, err := ob.client.RunQuery(query)
	if err != nil {
		return 0, fmt.Errorf("running query failed: %w", err)
	}

	ms := time.Duration(int64(value)) * time.Millisecond
	return ms, nil
}
