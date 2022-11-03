// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v0200Metrics, _ = os.ReadFile("testdata/v0.20.0/metrics.txt")
)

func Test_TestData(t *testing.T) {
	for name, data := range map[string][]byte{
		"v0200Metrics": v0200Metrics,
	} {
		assert.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.IsType(t, (*WMI)(nil), New())
}

func TestWMI_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"success if 'url' is set": {
			config: Config{
				HTTP: web.HTTP{Request: web.Request{URL: "http://127.0.0.1:9182/metrics"}}},
		},
		"fails on default config": {
			wantFail: true,
			config:   New().Config,
		},
		"fails if 'url' is unset": {
			wantFail: true,
			config:   Config{HTTP: web.HTTP{Request: web.Request{URL: ""}}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wmi := New()
			wmi.Config = test.config

			if test.wantFail {
				assert.False(t, wmi.Init())
			} else {
				assert.True(t, wmi.Init())
			}
		})
	}
}

func TestWMI_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func() (wmi *WMI, cleanup func())
		wantFail bool
	}{
		"success on valid response v0.20.0": {
			prepare: prepareWMIv0200,
		},
		"fails if endpoint returns invalid data": {
			wantFail: true,
			prepare:  prepareWMIReturnsInvalidData,
		},
		"fails on connection refused": {
			wantFail: true,
			prepare:  prepareWMIConnectionRefused,
		},
		"fails on 404 response": {
			wantFail: true,
			prepare:  prepareWMIResponse404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wmi, cleanup := test.prepare()
			defer cleanup()

			require.True(t, wmi.Init())

			if test.wantFail {
				assert.False(t, wmi.Check())
			} else {
				assert.True(t, wmi.Check())
			}
		})
	}
}

func TestWMI_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestWMI_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestWMI_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() (wmi *WMI, cleanup func())
		wantCollected map[string]int64
	}{
		"success on valid response v0.20.0": {
			prepare: prepareWMIv0200,
			wantCollected: map[string]int64{
				"collector_cpu_duration":                                        0,
				"collector_cpu_status_fail":                                     0,
				"collector_cpu_status_success":                                  1,
				"collector_logical_disk_duration":                               0,
				"collector_logical_disk_status_fail":                            0,
				"collector_logical_disk_status_success":                         1,
				"collector_logon_duration":                                      113,
				"collector_logon_status_fail":                                   0,
				"collector_logon_status_success":                                1,
				"collector_memory_duration":                                     0,
				"collector_memory_status_fail":                                  0,
				"collector_memory_status_success":                               1,
				"collector_net_duration":                                        0,
				"collector_net_status_fail":                                     0,
				"collector_net_status_success":                                  1,
				"collector_os_duration":                                         2,
				"collector_os_status_fail":                                      0,
				"collector_os_status_success":                                   1,
				"collector_process_duration":                                    115,
				"collector_process_status_fail":                                 0,
				"collector_process_status_success":                              1,
				"collector_service_duration":                                    101,
				"collector_service_status_fail":                                 0,
				"collector_service_status_success":                              1,
				"collector_system_duration":                                     0,
				"collector_system_status_fail":                                  0,
				"collector_system_status_success":                               1,
				"collector_tcp_duration":                                        0,
				"collector_tcp_status_fail":                                     0,
				"collector_tcp_status_success":                                  1,
				"cpu_core_0,0_cstate_c1":                                        320466854,
				"cpu_core_0,0_cstate_c2":                                        0,
				"cpu_core_0,0_cstate_c3":                                        0,
				"cpu_core_0,0_dpc_time":                                         134218,
				"cpu_core_0,0_dpcs":                                             9743800,
				"cpu_core_0,0_idle_time":                                        324911186,
				"cpu_core_0,0_interrupt_time":                                   154562,
				"cpu_core_0,0_interrupts":                                       310388662,
				"cpu_core_0,0_privileged_time":                                  2364218,
				"cpu_core_0,0_user_time":                                        2147342,
				"cpu_core_0,1_cstate_c1":                                        319056108,
				"cpu_core_0,1_cstate_c2":                                        0,
				"cpu_core_0,1_cstate_c3":                                        0,
				"cpu_core_0,1_dpc_time":                                         22186,
				"cpu_core_0,1_dpcs":                                             3301104,
				"cpu_core_0,1_idle_time":                                        318956250,
				"cpu_core_0,1_interrupt_time":                                   116186,
				"cpu_core_0,1_interrupts":                                       158651694,
				"cpu_core_0,1_privileged_time":                                  3602468,
				"cpu_core_0,1_user_time":                                        6864000,
				"cpu_core_0,2_cstate_c1":                                        319783446,
				"cpu_core_0,2_cstate_c2":                                        0,
				"cpu_core_0,2_cstate_c3":                                        0,
				"cpu_core_0,2_dpc_time":                                         32124,
				"cpu_core_0,2_dpcs":                                             4472938,
				"cpu_core_0,2_idle_time":                                        319696874,
				"cpu_core_0,2_interrupt_time":                                   107030,
				"cpu_core_0,2_interrupts":                                       134610838,
				"cpu_core_0,2_privileged_time":                                  3625092,
				"cpu_core_0,2_user_time":                                        6100500,
				"cpu_core_0,3_cstate_c1":                                        319088234,
				"cpu_core_0,3_cstate_c2":                                        0,
				"cpu_core_0,3_cstate_c3":                                        0,
				"cpu_core_0,3_dpc_time":                                         16280,
				"cpu_core_0,3_dpcs":                                             2370092,
				"cpu_core_0,3_idle_time":                                        319055092,
				"cpu_core_0,3_interrupt_time":                                   88968,
				"cpu_core_0,3_interrupts":                                       121533876,
				"cpu_core_0,3_privileged_time":                                  3521656,
				"cpu_core_0,3_user_time":                                        6845750,
				"cpu_dpc_time":                                                  204808,
				"cpu_idle_time":                                                 1282619402,
				"cpu_interrupt_time":                                            466746,
				"cpu_privileged_time":                                           13113434,
				"cpu_user_time":                                                 21957592,
				"logical_disk_C:_free_space":                                    43636490240,
				"logical_disk_C:_read_bytes_total":                              17676328448,
				"logical_disk_C:_read_latency":                                  97,
				"logical_disk_C:_reads_total":                                   350593,
				"logical_disk_C:_total_space":                                   67938287616,
				"logical_disk_C:_used_space":                                    24301797376,
				"logical_disk_C:_write_bytes_total":                             9135282688,
				"logical_disk_C:_write_latency":                                 123,
				"logical_disk_C:_writes_total":                                  450705,
				"logon_type_batch_sessions":                                     0,
				"logon_type_cached_interactive_sessions":                        0,
				"logon_type_cached_remote_interactive_sessions":                 0,
				"logon_type_cached_unlock_sessions":                             0,
				"logon_type_interactive_sessions":                               2,
				"logon_type_network_clear_text_sessions":                        0,
				"logon_type_network_sessions":                                   0,
				"logon_type_new_credentials_sessions":                           0,
				"logon_type_proxy_sessions":                                     0,
				"logon_type_remote_interactive_sessions":                        0,
				"logon_type_service_sessions":                                   0,
				"logon_type_system_sessions":                                    0,
				"logon_type_unlock_sessions":                                    0,
				"memory_available_bytes":                                        1379942400,
				"memory_cache_bytes":                                            170774528,
				"memory_cache_bytes_peak":                                       208621568,
				"memory_cache_faults_total":                                     8009603,
				"memory_cache_total":                                            1392185344,
				"memory_commit_limit":                                           5733113856,
				"memory_committed_bytes":                                        3447439360,
				"memory_demand_zero_faults_total":                               102505136,
				"memory_free_and_zero_page_list_bytes":                          20410368,
				"memory_free_system_page_table_entries":                         16722559,
				"memory_modified_page_list_bytes":                               32653312,
				"memory_not_committed_bytes":                                    2285674496,
				"memory_page_faults_total":                                      119093924,
				"memory_pool_nonpaged_allocs_total":                             0,
				"memory_pool_nonpaged_bytes_total":                              126865408,
				"memory_pool_paged_allocs_total":                                0,
				"memory_pool_paged_bytes":                                       303906816,
				"memory_pool_paged_resident_bytes":                              294293504,
				"memory_standby_cache_core_bytes":                               107376640,
				"memory_standby_cache_normal_priority_bytes":                    1019121664,
				"memory_standby_cache_reserve_bytes":                            233033728,
				"memory_standby_cache_total":                                    1359532032,
				"memory_swap_page_operations_total":                             4956175,
				"memory_swap_page_reads_total":                                  402087,
				"memory_swap_page_writes_total":                                 7012,
				"memory_swap_pages_read_total":                                  4643279,
				"memory_swap_pages_written_total":                               312896,
				"memory_system_cache_resident_bytes":                            170774528,
				"memory_system_code_resident_bytes":                             17100800,
				"memory_system_code_total_bytes":                                8192,
				"memory_system_driver_resident_bytes":                           46092288,
				"memory_system_driver_total_bytes":                              18731008,
				"memory_transition_faults_total":                                27183527,
				"memory_transition_pages_repurposed_total":                      2856471,
				"memory_used_bytes":                                             2876776448,
				"memory_write_copies_total":                                     1194039,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_bytes_received": 76581511712,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_bytes_sent":     16422331008,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_bytes_total":    93003842720,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_outbound_discarded": 0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_outbound_errors":    0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_received_discarded": 0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_received_errors":    0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_received_total":     8241738,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_received_unknown":   0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_sent_total":         2664932,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_total":              10906670,
				"os_paging_free_bytes":                     1414107136,
				"os_paging_limit_bytes":                    1476395008,
				"os_paging_used_bytes":                     62287872,
				"os_physical_memory_free_bytes":            1379946496,
				"os_process_memory_limit_bytes":            140737488224256,
				"os_processes":                             152,
				"os_processes_limit":                       4294967295,
				"os_time":                                  1667508748,
				"os_users":                                 2,
				"os_virtual_memory_bytes":                  5733113856,
				"os_virtual_memory_free_bytes":             2285674496,
				"os_visible_memory_bytes":                  4256718848,
				"os_visible_memory_used_bytes":             2876772352,
				"process_msedge_cpu_time":                  3839786,
				"process_msedge_handles":                   11558,
				"process_msedge_io_bytes":                  7956454756,
				"process_msedge_io_operations":             33477284,
				"process_msedge_page_faults":               10711882,
				"process_msedge_page_file_bytes":           1363206144,
				"process_msedge_threads":                   426,
				"process_msedge_working_set_private_bytes": 922689536,
				"service_dhcp_state_continue_pending":      0,
				"service_dhcp_state_pause_pending":         0,
				"service_dhcp_state_paused":                0,
				"service_dhcp_state_running":               1,
				"service_dhcp_state_start_pending":         0,
				"service_dhcp_state_stop_pending":          0,
				"service_dhcp_state_stopped":               0,
				"service_dhcp_state_unknown":               0,
				"service_dhcp_status_degraded":             0,
				"service_dhcp_status_error":                0,
				"service_dhcp_status_lost_comm":            0,
				"service_dhcp_status_no_contact":           0,
				"service_dhcp_status_nonrecover":           0,
				"service_dhcp_status_ok":                   1,
				"service_dhcp_status_pred_fail":            0,
				"service_dhcp_status_service":              0,
				"service_dhcp_status_starting":             0,
				"service_dhcp_status_stopping":             0,
				"service_dhcp_status_stressed":             0,
				"service_dhcp_status_unknown":              0,
				"system_calls_total":                       1886567439,
				"system_context_switches_total":            486550330,
				"system_exception_dispatches_total":        160348,
				"system_processor_queue_length":            0,
				"system_threads":                           1559,
				"system_up_time":                           170968,
				"tcp_ipv4_conns_active":                    4301,
				"tcp_ipv4_conns_established":               7,
				"tcp_ipv4_conns_failures":                  137,
				"tcp_ipv4_conns_passive":                   501,
				"tcp_ipv4_conns_resets":                    1282,
				"tcp_ipv4_segments_received":               676388,
				"tcp_ipv4_segments_retransmitted":          2120,
				"tcp_ipv4_segments_sent":                   871379,
				"tcp_ipv6_conns_active":                    214,
				"tcp_ipv6_conns_established":               0,
				"tcp_ipv6_conns_failures":                  214,
				"tcp_ipv6_conns_passive":                   0,
				"tcp_ipv6_conns_resets":                    0,
				"tcp_ipv6_segments_received":               1284,
				"tcp_ipv6_segments_retransmitted":          428,
				"tcp_ipv6_segments_sent":                   856,
			},
		},
		"fails if endpoint returns invalid data": {
			prepare: prepareWMIReturnsInvalidData,
		},
		"fails on connection refused": {
			prepare: prepareWMIConnectionRefused,
		},
		"fails on 404 response": {
			prepare: prepareWMIResponse404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wmi, cleanup := test.prepare()
			defer cleanup()

			require.True(t, wmi.Init())

			mx := wmi.Collect()

			if mx != nil && test.wantCollected != nil {
				mx["system_up_time"] = test.wantCollected["system_up_time"]
			}

			assert.Equal(t, test.wantCollected, mx)
			if len(test.wantCollected) > 0 {
				testCharts(t, wmi, mx)
			}
		})
	}
}
func testCharts(t *testing.T, wmi *WMI, mx map[string]int64) {
	ensureChartsDimsCreated(t, wmi)
	ensureCollectedHasAllChartsDimsVarsIDs(t, wmi, mx)
}

func ensureChartsDimsCreated(t *testing.T, w *WMI) {
	for _, chart := range cpuCharts {
		if w.cache.collection[collectorCPU] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range memCharts {
		if w.cache.collection[collectorMemory] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range tcpCharts {
		if w.cache.collection[collectorTCP] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range osCharts {
		if w.cache.collection[collectorOS] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range systemCharts {
		if w.cache.collection[collectorSystem] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range logonCharts {
		if w.cache.collection[collectorLogon] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range processesCharts {
		if w.cache.collection[collectorProcess] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}

	for core := range w.cache.cores {
		for _, chart := range cpuCoreChartsTmpl {
			id := fmt.Sprintf(chart.ID, core)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' core", id, core)
		}
	}
	for disk := range w.cache.volumes {
		for _, chart := range diskChartsTmpl {
			id := fmt.Sprintf(chart.ID, disk)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' disk", id, disk)
		}
	}
	for nic := range w.cache.nics {
		for _, chart := range nicChartsTmpl {
			id := fmt.Sprintf(chart.ID, nic)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' nic", id, nic)
		}
	}
	for zone := range w.cache.thermalZones {
		for _, chart := range thermalzoneChartsTmpl {
			id := fmt.Sprintf(chart.ID, zone)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' thermalzone", id, zone)
		}
	}
	for svc := range w.cache.services {
		for _, chart := range serviceChartsTmpl {
			id := fmt.Sprintf(chart.ID, svc)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' service", id, svc)
		}
	}
	for name := range w.cache.collectors {
		for _, chart := range collectorChartsTmpl {
			id := fmt.Sprintf(chart.ID, name)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' collector", id, name)
		}
	}

	for _, chart := range processesCharts {
		if chart = w.Charts().Get(chart.ID); chart == nil {
			continue
		}
		for proc := range w.cache.processes {
			var found bool
			for _, dim := range chart.Dims {
				if found = strings.HasPrefix(dim.ID, "process_"+proc); found {
					break
				}
			}
			assert.Truef(t, found, "chart '%s' has not dim for '%s' process", chart.ID, proc)
		}
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, w *WMI, mx map[string]int64) {
	for _, chart := range *w.Charts() {
		for _, dim := range chart.Dims {
			_, ok := mx[dim.ID]
			assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := mx[v.ID]
			assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
}

func prepareWMIv0200() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(v0200Metrics)
		}))

	wmi = New()
	wmi.URL = ts.URL
	return wmi, ts.Close
}

func prepareWMIReturnsInvalidData() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))

	wmi = New()
	wmi.URL = ts.URL
	return wmi, ts.Close
}

func prepareWMIConnectionRefused() (wmi *WMI, cleanup func()) {
	wmi = New()
	wmi.URL = "http://127.0.0.1:38001"
	return wmi, func() {}
}

func prepareWMIResponse404() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

	wmi = New()
	wmi.URL = ts.URL
	return wmi, ts.Close
}
