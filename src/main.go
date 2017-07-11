/*
 * ip as-path regex builder
 * Jörg Kost, joerg.kost@gmx.com, jk@ip-clear.de
 *
 */

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sort"
	"./numberRange"
	"os"
)

var NicToASN map[string][]string
var fmt_asPathACL = "ip as-path access-list %s %s %s$\n"

/* Command line parameters */
var flagNicRegion string
var flagPermitOrDeny int
var flagAclTitle string
var flagSummaryOnly bool
var flagNICsParsed []string
var flagFilename string
var PermitOrDenyArr [2]string = [2]string{"deny", "permit"}

var asnList []string = []string{
	"http://www.iana.org/assignments/as-numbers/as-numbers-2.csv",
	"http://www.iana.org/assignments/as-numbers/as-numbers-1.csv"}

func init() {
	NicToASN = make(map[string][]string)
	flag.StringVar(&flagNicRegion, "region", "", "Comma separated list with region for generated prefix")
	flag.StringVar(&flagAclTitle, "acltitle", "region-summary", "Title for generated as-path list")
	flag.BoolVar(&flagSummaryOnly, "summary", false, "Print summary of downloaded lists only")
	flag.IntVar(&flagPermitOrDeny, "permitOrDeny", 1, "Deny = 0, Permit = 1")
	flag.StringVar(&flagFilename, "filename", "", "Output file, else stdout")
	flag.Parse()
	flagNICsParsed = strings.Split(flagNicRegion, ",")
}

func main() {
	// log.Println("Building ASN-prefixlists per NIC")
	var wg sync.WaitGroup
	var mt sync.Mutex

	for _, v := range asnList {
		wg.Add(1)
		go func(asnURI string) {
			defer wg.Done()
			// log.Printf("downloading %s\n", asnURI)

			resp, err := http.Get(asnURI)
			if err != nil {
				log.Fatal("Cant open as numbers from IANA")
			}

			mt.Lock()
			map_asn_to_nic(resp.Body)
			mt.Unlock()
			resp.Body.Close()
		}(v)
	}

	wg.Wait()
	if flagSummaryOnly == true || len(flagNICsParsed) == 0 {
		printSummary()
	} else {
		if flagFilename == "" {
			generatePrefixList(os.Stdout)
		} else {
			outputStream, err := os.Create(flagFilename)
			if err != nil {
				log.Fatal(err)
			}
			generatePrefixList(outputStream)
			outputStream.Close()
		}
	}
}

func printSummary() {
	for k := range NicToASN {
		log.Printf("%s [%d table entries]\n", k, len(NicToASN[k]))
	}
}

func generatePrefixList(outputStream io.Writer ) {
	var prefixLists []string
	for _, nic := range flagNICsParsed {
		for _, v := range NicToASN[nic] {
			if strings.Contains(v, "-") {
				rangeSplit := strings.Split(v, "-")
				start, err := strconv.Atoi(rangeSplit[0])
				if err != nil {
					panic(err)
				}
				end, err := strconv.Atoi(rangeSplit[1])
				if err != nil {
					panic(err)
				}
				if start == end {
					prefixLists = append(prefixLists, fmt.Sprintf(fmt_asPathACL, flagAclTitle, PermitOrDenyArr[flagPermitOrDeny],
						"_" + strconv.Itoa(start)))
				} else {
					prefixLists = append(prefixLists, fmt.Sprintf(fmt_asPathACL, flagAclTitle, PermitOrDenyArr[flagPermitOrDeny],
						numberRange.GetRegex(start, end)))
				}
			} else {
				prefixLists = append(prefixLists, fmt.Sprintf(fmt_asPathACL, flagAclTitle, PermitOrDenyArr[flagPermitOrDeny], "_" + v))
			}
		}
	}
	sort.Slice(prefixLists,
		func(i,j int) bool {
			return prefixLists[i] < prefixLists[j]
		});

	for _, v := range prefixLists {
		fmt.Fprint(outputStream, v)
	}
}

func map_asn_to_nic(ASN io.ReadCloser) {
	r := csv.NewReader(ASN)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasPrefix(record[1], "Assigned by ") {
			nic := strings.TrimPrefix(record[1], "Assigned by ")
			NicToASN[nic] = append(NicToASN[nic], record[0])
		} else if strings.HasPrefix(record[1], "Reserved") || strings.HasPrefix(record[1], "Unallocated") {
			NicToASN["bogons"] = append(NicToASN["bogons"], record[0])
 		}
	}
}



