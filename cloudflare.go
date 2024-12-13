package main

import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "text/tabwriter"
)

func main() {
    mode := flag.String("mode", "", "Mode for the firewall rule")
    ip := flag.String("ip", "", "IP address for the firewall rule")
    notes := flag.String("notes", "", "Notes for the firewall rule")
    del := flag.Bool("del", false, "Delete a firewall rule")
    flag.Parse()

    account := "test@domain.com"
    key := "123"

    client := &http.Client{}

    if *mode == "" && *ip == "" && !*del {
        // Perform GET request to check existing rules
        url := "https://api.cloudflare.com/client/v4/user/firewall/access_rules/rules"
        method := "GET"

        req, err := http.NewRequest(method, url, nil)
        if err != nil {
            fmt.Println("Error creating GET request:", err)
            return
        }

        req.Header.Add("Content-Type", "application/json")
        req.Header.Add("X-Auth-Email", account)
        req.Header.Add("X-Auth-Key", key)

        res, err := client.Do(req)
        if err != nil {
            fmt.Println("Error executing GET request:", err)
            return
        }
        defer res.Body.Close()

        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
            fmt.Println("Error reading GET response body:", err)
            return
        }

        var result map[string]interface{}
        if err := json.Unmarshal(body, &result); err != nil {
            fmt.Println("Error parsing GET response JSON:", err)
            return
        }

        if success, ok := result["success"].(bool); ok && !success {
            fmt.Println("API returned an error:", result["errors"])
            return
        }

        // Print the response in a formatted table
        writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
        fmt.Fprintln(writer, "ID\tIP\tMode\tDate\tNotes")
        fmt.Fprintln(writer, "================================\t===============\t====\t====\t=====")
        for _, rule := range result["result"].([]interface{}) {
            ruleMap := rule.(map[string]interface{})
            id := ruleMap["id"]
            ip := ruleMap["configuration"].(map[string]interface{})["value"]
            mode := ruleMap["mode"]
            date := ruleMap["created_on"]
            notes := ruleMap["notes"]
            fmt.Fprintf(writer, "%v\t%v\t%v\t%v\t%v\n", id, ip, mode, date, notes)
        }
        writer.Flush()
    } else if *mode != "" && *ip != "" && !*del {
        // Perform GET request to check existing rules
        url := "https://api.cloudflare.com/client/v4/user/firewall/access_rules/rules"
        method := "GET"

        req, err := http.NewRequest(method, url, nil)
        if err != nil {
            fmt.Println("Error creating GET request:", err)
            return
        }

        req.Header.Add("Content-Type", "application/json")
        req.Header.Add("X-Auth-Email", account)
        req.Header.Add("X-Auth-Key", key)

        res, err := client.Do(req)
        if err != nil {
            fmt.Println("Error executing GET request:", err)
            return
        }
        defer res.Body.Close()

        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
            fmt.Println("Error reading GET response body:", err)
            return
        }

        var result map[string]interface{}
        if err := json.Unmarshal(body, &result); err != nil {
            fmt.Println("Error parsing GET response JSON:", err)
            return
        }

        if success, ok := result["success"].(bool); ok && !success {
            fmt.Println("API returned an error:", result["errors"])
            return
        }

        ipExists := false
        if success, ok := result["success"].(bool); ok && success {
            if rules, ok := result["result"].([]interface{}); ok {
                for _, rule := range rules {
                    if ruleMap, ok := rule.(map[string]interface{}); ok {
                        if ruleMap["configuration"].(map[string]interface{})["value"] == *ip {
                            ipExists = true
                            break
                        }
                    }
                }
            }
        } else {
            fmt.Println("Failed to retrieve rules")
            return
        }

        if ipExists {
            fmt.Println("IP address already exists in the rules. Skipping POST request.")
            return
        }

        // Perform POST request
        url = "https://api.cloudflare.com/client/v4/user/firewall/access_rules/rules"
        method = "POST"

        data := fmt.Sprintf(`{
            "configuration": {
                "target": "ip",
                "value": "%s"
            },
            "mode": "%s",
            "notes": "%s"
        }`, *ip, *mode, *notes)

        req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
        if err != nil {
            fmt.Println("Error creating POST request:", err)
            return
        }

        req.Header.Add("Content-Type", "application/json")
        req.Header.Add("X-Auth-Email", account)
        req.Header.Add("X-Auth-Key", key)

        res, err = client.Do(req)
        if err != nil {
            fmt.Println("Error executing POST request:", err)
            return
        }
        defer res.Body.Close()

        fmt.Println("Response Status:", res.Status)
    } else if *del && *mode != "" && *ip != "" {
        // Perform GET request to find the rule ID
        url := "https://api.cloudflare.com/client/v4/user/firewall/access_rules/rules"
        method := "GET"

        req, err := http.NewRequest(method, url, nil)
        if err != nil {
            fmt.Println("Error creating GET request:", err)
            return
        }

        req.Header.Add("Content-Type", "application/json")
        req.Header.Add("X-Auth-Email", account)
        req.Header.Add("X-Auth-Key", key)

        res, err := client.Do(req)
        if err != nil {
            fmt.Println("Error executing GET request:", err)
            return
        }
        defer res.Body.Close()

        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
            fmt.Println("Error reading GET response body:", err)
            return
        }

        var result map[string]interface{}
        if err := json.Unmarshal(body, &result); err != nil {
            fmt.Println("Error parsing GET response JSON:", err)
            return
        }

        if success, ok := result["success"].(bool); ok && !success {
            fmt.Println("API returned an error:", result["errors"])
            return
        }

        var ruleID string
        if success, ok := result["success"].(bool); ok && success {
            if rules, ok := result["result"].([]interface{}); ok {
                for _, rule := range rules {
                    if ruleMap, ok := rule.(map[string]interface{}); ok {
                        if ruleMap["configuration"].(map[string]interface{})["value"] == *ip && ruleMap["mode"] == *mode {
                            ruleID = ruleMap["id"].(string)
                            break
                        }
                    }
                }
            }
        } else {
            fmt.Println("Failed to retrieve rules")
            return
        }

        if ruleID == "" {
            fmt.Println("No matching rule found for deletion")
            return
        }

        // Perform DELETE request
        url = fmt.Sprintf("https://api.cloudflare.com/client/v4/user/firewall/access_rules/rules/%s", ruleID)
        method = "DELETE"

        req, err = http.NewRequest(method, url, nil)
        if err != nil {
            fmt.Println("Error creating DELETE request:", err)
            return
        }

        req.Header.Add("Content-Type", "application/json")
        req.Header.Add("X-Auth-Email", account)
        req.Header.Add("X-Auth-Key", key)

        res, err = client.Do(req)
        if err != nil {
            fmt.Println("Error executing DELETE request:", err)
            return
        }
        defer res.Body.Close()

        fmt.Println("Response Status:", res.Status)
    } else {
        fmt.Println("Invalid input. Use either GET request without parameters, POST request with --mode, --ip, and --notes parameters, or DELETE request with --del, --mode, and --ip parameters.")
        flag.Usage()
        os.Exit(1)
    }
}
