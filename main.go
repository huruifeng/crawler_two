package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/semaphore"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//LIN – Lincoln Service Center (Nebraska Service Center)
//EAC – Eastern Adjudication Center (Vermont Service Center)
//IOE – USCIS Electronic Immigration System (ELIS)
//WAC – Western Adjudication Center (California Service Center)
//MSC – Missouri Service Center (National Benefits Center)
//NBC – National Benefits Center
//SRC – Southern Regional Center (Texas Service Center)
//YSC - Potomac Service Center
//

var CENTER_NAMES = []string{
	"LIN",
	"MSC",
	"SRC",
	"WAC",
	"EAC",
	"YSC",
	//"IOE",
}

var MONTHS = map[string]bool{
	"January":   true,
	"February":  true,
	"March":     true,
	"April":     true,
	"May":       true,
	"June":      true,
	"July":      true,
	"August":    true,
	"September": true,
	"October":   true,
	"November":  true,
	"December":  true,
}

var FORM_TYPES = []string{
	"I-129CW",
	"I-129F",
	"I-290B",
	"I-600A",
	"I-601A",
	"I-751A",
	"I-765V",
	"I-485J",
	"I-800A",
	"I-821D",
	"I-90",
	"I-102",
	"I-129",
	"I-130",
	"I-131",
	"I-140",
	"I-212",
	"I-360",
	"I-485",
	"I-526",
	"I-539",
	"I-600",
	"I-601",
	"I-612",
	"I-730",
	"I-751",
	"I-765",
	"I-800",
	"I-817",
	"I-821",
	"I-824",
	"I-829",
	"I-914",
	"I-918",
	"I-924",
	"I-929",
	"EOIR-29",
	"G-28",
}

var FINAL_STATUS = map[string]bool{
	"CASE STATUS":          true,
	"Appeal Was Dismissed": true,
	"Card Destroyed":       true,
	"Card Was Destroyed":   true,
	"Document Destroyed":   true,
	"Appeal Was Approved":  true,
	"Case Was Approved":    true,
	"Case Was Denied":      true,
	"Special Immigrant Juvenile Deferred Action Granted":                             true,
	"Case Approval Was Affirmed":                                                     true,
	"Case Approval Was Certified By USCIS":                                           true,
	"Case Approval Was Reaffirmed And Mailed Back To Department Of State":            true,
	"Case Was Approved And A Decision Notice Was Sent":                               true,
	"Case Was Approved And My Decision Was Emailed":                                  true,
	"Case Was Approved And USCIS Notified The U.S. Consulate or Port of Entry":       true,
	"Request For A Duplicate Card Was Approved":                                      true,
	"Special Immigrant Juvenile Form I-360 Approved With Deferred Action":            true,
	"Card Was Delivered To Me By The Post Office":                                    true,
	"Card Was Determined As Undeliverable By The Post Office":                        true,
	"Card Was Mailed To Me":                                                          true,
	"Card Was Picked Up By The United States Postal Service":                         true,
	"Document Was Mailed":                                                            true,
	"Document Was Mailed But Not Returned To USCIS":                                  true,
	"Document Was Mailed To Me":                                                      true,
	"Document Was Personally Delivered To Me":                                        true,
	"Travel Document Was Mailed":                                                     true,
	"Advance Parole Document Was Produced":                                           true,
	"New Card Is Being Produced":                                                     true,
	"Permanent Resident Card Is Being Produced":                                      true,
	"Reentry Permit Was Produced":                                                    true,
	"Refugee Travel Document Was Produced":                                           true,
	"Document Is Being Held For 180 Days":                                            true,
	"Document Was Returned As Undeliverable":                                         true,
	"Document Was Returned To USCIS":                                                 true,
	"Document Was Returned To USCIS And Is Being Held":                               true,
	"Special Immigrant Juvenile Deferred Action Not Granted":                         true,
	"Certified Approval Of My Case Was Reversed by The Appellate Authority":          true,
	"Certified Denial Of My Case Was Affirmed By Appellate Authority":                true,
	"Case Was Denied And My Decision Notice Mailed":                                  true,
	"Denial Was Upheld by Court":                                                     true,
	"Case Rejected Because I Sent An Incorrect Fee":                                  true,
	"Case Rejected Because The Version Of The Form I Sent Is No Longer Accepted":     true,
	"Case Rejected For Form Not Signed And Incorrect Form Version":                   true,
	"Case Rejected For Incorrect Fee And Form Not Signed":                            true,
	"Case Rejected For Incorrect Fee And Incorrect Form Version":                     true,
	"Case Rejected For Incorrect Fee And Payment Not Signed":                         true,
	"Case Rejected For Incorrect Fee, Payment Not Signed And Incorrect Form Version": true,
	"Case Was Rejected Because I Did Not Sign My Form":                               true,
	"Case Was Rejected Because It Was Improperly Filed":                              true,
	"Case Was Rejected Because My Check Or Money Order Is Not Signed":                true,
	"CNMI Semiannual Report Rejected":                                                true,
	"Form G-28 Was Rejected Because It Was Improperly Filed":                         true,
	"Petition/Application Was Rejected For Insufficient Funds":                       true,
	"Case Approval Was Revoked":                                                      true,
	"Special Immigrant Juvenile Deferred Action Terminated":                          true,
	"Document Was Destroyed After USCIS Held It For 180 Days":                        true,
	"Document Was Destroyed And Letter Was Received":                                 true,
	"Travel Document Was Destroyed":                                                  true,
	"Travel Document Was Destroyed After USCIS Held It For 180 Days":                 true,
	"Student Employment Authorization Document Automatically Terminated":             true,
	"Case Is No Longer On Hold Because Of Pending Litigation":                        true,
	"Termination Of Litigation Notice Was Mailed":                                    true,
	"Card Is Being Returned to USCIS by Post Office":                                 true,
	"Card Was Returned To USCIS":                                                     true,
	"Notice Was Returned To USCIS Because The Post Office Could Not Deliver It":      true,
	"Travel Document Was Returned to USCIS And Will Be Held For 180 Days":            true,
	"Case Was Automatically Revoked":                                                 true,
	"Revocation Notice Was Sent":                                                     true,
	"Appeal Was Terminated and A Notice Was Mailed To Me":                            true,
	"Case Closed Benefit Received By Other Means":                                    true,
	"Petition Business Terminated/Over 180 Days/Not Automatically Revoked":           true,
	"Termination Notice Sent":                                                        true,
	"Petition Withdrawn/Over 180 Days/Not Automatically Revoked":                     true,
}

type result struct {
	case_id string
	status  string
	form    string
	date    string
}

//var day_case_count_mutex sync.Mutex
//var day_case_count = make(map[int]int)

var case_status_store_mutex sync.Mutex
var case_status_store = make(map[string][]string)

var case_final_store_mutex sync.Mutex
var case_final_store_temp = make(map[string][]string)
var case_final_store = make(map[string][]string)

var sem = semaphore.NewWeighted(1000)

var done_n = 0
var try_n = 0

func Split(r rune) bool {
	return r == ',' || r == ' '
}

func writeF(path string, content []byte) {
	err := os.WriteFile(path, content, 0666)
	if err != nil {
		fmt.Println("Write error! ", err.Error())
	}
}

func get(form url.Values, retry int, try_times int) result {
	case_id := form.Get("appReceiptNum")
	if try_times > 0 {
		if retry > try_times {
			fmt.Printf("%s:%s-%d", case_id, "try faild", try_times)
			return result{case_id, "try_faild", "", ""}
		}
	}

	sem.Acquire(context.Background(), 1)
	res, err1 := http.PostForm("https://egov.uscis.gov/casestatus/mycasestatus.do", form)
	sem.Release(1)

	defer func() {
		if err1 == nil {
			res.Body.Close()
		}
	}()
	if err1 != nil {
		//fmt.Println("error 1! " + err1.Error())
		//fmt.Printf("Retry %d %s\n", retry+1, form)
		return get(form, retry+1, try_times)
	}

	doc, err2 := goquery.NewDocumentFromReader(res.Body)
	if err2 != nil {
		//fmt.Println("error 2! " + err2.Error())
		//fmt.Printf("Retry %d %s\n", retry+1, form)
		return get(form, retry+1, try_times)
	}

	body := doc.Find(".rows").First()
	status_h := body.Find("h1").Text()
	status_p := body.Find("p").Text()
	status_p_s := strings.FieldsFunc(status_p, Split)
	date_x := ""
	form_x := ""
	for i, w := range status_p_s {
		if MONTHS[w] {
			date_x = status_p_s[i] + " " + status_p_s[i+1] + ", " + status_p_s[i+2]
		} else if w == "Form" {
			form_x = status_p_s[i+1]
			break
		}
	}
	if form_x == "" {
		for _, form_i_x := range FORM_TYPES {
			if strings.Contains(doc.Text(), form_i_x) {
				form_x = form_i_x
				break
			}
		}
	}

	if status_h != "" {
		return result{case_id, status_h, form_x, date_x}
	} else {
		return result{case_id, "invalid_num", "", ""}
	}
}

func buildURL(center string, two_digit_yr int, day int, code int, case_serial_numbers int, format string) url.Values {
	if format == "SC" {
		res := url.Values{"appReceiptNum": {fmt.Sprintf("%s%d%03d%d%04d", center, two_digit_yr, day, code, case_serial_numbers)}}
		return res
	} else {
		res := url.Values{"appReceiptNum": {fmt.Sprintf("%s%d%d%03d%04d", center, two_digit_yr, code, day, case_serial_numbers)}}
		return res
	}
}

func crawlerAsync(center string, two_digit_yr int, day int, code int, case_serial_numbers int, format string, c chan result) {
	c <- crawler(center, two_digit_yr, day, code, case_serial_numbers, format, try_n)
}

func crawler(center string, two_digit_yr int, day int, code int, case_serial_numbers int, format string, try_times int) result {
	url_x := buildURL(center, two_digit_yr, day, code, case_serial_numbers, format)
	res := get(url_x, 0, try_times)

	//if res.status != "invalid_num" {
	//	//fmt.Printf("%s: %s %s %s\n", url_x.Get("appReceiptNum"), res.date, res.form, res.status)
	//
	//	if !FINAL_STATUS[res.status] {
	//		case_status_store_mutex.Lock()
	//		case_status_store[res.case_id] = []string{res.form, res.date, res.status}
	//		case_status_store_mutex.Unlock()
	//	}
	//}
	return res
}

func getLastCaseNumber(center string, two_digit_yr int, day int, code int, format string) int {
	low := 1
	high := 1
	invalid_limit := 10

	i := 0
	for high < 10000 {
		for i = 0; i < invalid_limit; i++ {
			if crawler(center, two_digit_yr, day, code, high+i-1, format, 0).status != "invalid_num" {
				high *= 2
				break
			}
		}
		if i == invalid_limit {
			break
		}
	}

	for low < high {
		mid := (low + high) / 2
		for i = 0; i < invalid_limit; i++ {
			if crawler(center, two_digit_yr, day, code, mid+i, format, 0).status != "invalid_num" {
				low = mid + 1
				break
			}
		}

		if i == invalid_limit {
			high = mid
		}
	}

	return low - 1
}

func merge_final_case(cur_cases, new_cases map[string][]string) {
	for id_key, s_val := range new_cases {
		//if cur_cases[id_key] == nil {
		//	cur_cases[id_key] = s_val
		//} else {
		//	cur_cases[id_key] = s_val
		//}
		cur_cases[id_key] = s_val
	}
}

func all(center string, two_digit_yr int, day int, code int, format string, report_c chan int) {
	defer func() { report_c <- 0 }()

	last := getLastCaseNumber(center, two_digit_yr, day, code, format)
	fmt.Printf("loading %s-%s-%d at day %3d: %d\n", center, format, two_digit_yr, day, last)

	c := make(chan result)
	case_id := ""
	true_last := 0
	for i := 0; i <= last; i++ {
		if format == "SC" {
			case_id = fmt.Sprintf("%s%d%03d%d%04d", center, two_digit_yr, day, code, i)
		} else {
			case_id = fmt.Sprintf("%s%d%d%03d%04d", center, two_digit_yr, code, day, i)
		}
		_, has := case_final_store_temp[case_id]
		if !has {
			true_last += 1
			go crawlerAsync(center, two_digit_yr, day, code, i, format, c)
		}
	}

	new_final_status_case := make(map[string][]string)
	for i := 0; i < true_last; i++ {

		cur := <-c
		if cur.status == "invalid_num" {
			continue
		}
		//fmt.Sprintf("%s:%s|%s|%s", cur.case_id, cur.form, cur.date, cur.status)
		if FINAL_STATUS[cur.status] {
			new_final_status_case[cur.case_id] = []string{cur.form, cur.date, cur.status}
		} else {
			case_status_store_mutex.Lock()
			case_status_store[cur.case_id] = []string{cur.form, cur.date, cur.status}
			case_status_store_mutex.Unlock()
		}

	}

	case_final_store_mutex.Lock()
	done_n += 1
	merge_final_case(case_final_store, new_final_status_case)
	case_final_store_mutex.Unlock()

	//day_case_count_mutex.Lock()
	//day_case_count[day] = last
	//day_case_count_mutex.Unlock()

	fmt.Printf("Done %s-%s-%d at day %3d: %d =>> %d / 366\n", center, format, two_digit_yr, day, last, done_n)
}

func main() {
	args := os.Args

	fmt.Println(time.Now())
	center := args[1]
	fiscal_year, _ := strconv.Atoi(args[2])
	format := args[3]

	try_n, _ = strconv.Atoi(args[4])

	fmt.Printf("Run %s-%s-%d, Try = %d\n", center, format, fiscal_year, try_n)

	year_days := 365
	done_n = 0

	dir, _ := os.Getwd()

	// Load the final records of the center and fiscal year, to avoid visiting them again
	// final records are the cases with the statues defined in the FINAL_STATUS
	case_final_store_file := fmt.Sprintf("%s/saved_data/%s_%s_%d_case_final.json", dir, center, format, fiscal_year)
	jsonFile, err := os.ReadFile(case_final_store_file)
	if err != nil {
		fmt.Println("Read error! ", err.Error())
	} else {
		json.Unmarshal([]byte(jsonFile), &case_final_store)
		json.Unmarshal([]byte(jsonFile), &case_final_store_temp)
	}

	// Start the data retrieval
	if format == "LB" {
		report_c_lb := make(chan int)
		for day := 0; day <= year_days; day++ {
			go all(center, fiscal_year, day, 9, "LB", report_c_lb)
		}
		for i := 0; i <= year_days; i++ {
			<-report_c_lb
		}
	} else if format == "SC" {
		report_c_sc := make(chan int)
		for day := 0; day <= year_days; day++ {
			go all(center, fiscal_year, day, 5, "SC", report_c_sc)
		}
		for i := 0; i <= year_days; i++ {
			<-report_c_sc
		}
	} else if format == "ioe" {

	}

	// Save case status
	fmt.Println("Saving data...")
	case_status_save_path := fmt.Sprintf("%s/saved_data/%s_%s_%d.json", dir, center, format, fiscal_year)
	b_status, _ := json.MarshalIndent(case_status_store, "", "  ")
	writeF(case_status_save_path, b_status)

	// Save case with final status
	b_final, _ := json.MarshalIndent(case_final_store, "", "  ")
	writeF(case_final_store_file, b_final)
	fmt.Println("Saving data...Done!")
	//total := 0
	//for _, e := range day_case_count {
	//	total += e
	//
	//}
	//fmt.Println("Total:", total)

	fmt.Println(time.Now())
}
