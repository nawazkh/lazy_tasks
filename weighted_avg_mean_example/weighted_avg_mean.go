package main

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

type Coordinate struct {
	TotalCores      int
	WeightedAverage float64
	TotalCustomers  int
	ActionCount     int
}

var (
	globalCounter                     = 0
	moneyPerCorePerDay                = 10
	totalCustomers                    int
	maxNumberOfCoresThatCanBeConsumed int
	customersList                     = []string{}
	dailyCoreConsumptionPerCustomer   = map[string]int{}
	myPlot                            = []Coordinate{}
)

func main() {
	// start with random number of customers
	// - ask user for random number of customers
	// - display average total consumption
	// - ask user to input a new customer and total core consumtion
	// - display new core count.

	fmt.Println("Enter initial number of customers to consider")
	fmt.Scanf("%d", &totalCustomers)
	fmt.Println("")

	fmt.Println("Enter max number of cores that can be consumed by a customer")
	fmt.Scanf("%d", &maxNumberOfCoresThatCanBeConsumed)
	fmt.Println("")

	// make random customer list
	customersList = makeCustomerSortedList(totalCustomers)

	// not so random now
	rand.Seed(10)

	// Create a list for totalCustomers with random daily core count per customer
	for _, customerName := range customersList {
		dailyCoreConsumptionPerCustomer[customerName] = rand.Intn(maxNumberOfCoresThatCanBeConsumed) + 1
	}

	// Initial weighted mean to work with
	generateWeightedMean()

	var choice string
	fmt.Println("Would you like to add/delete customers? y/n")
	fmt.Scanf("%s", &choice)
	fmt.Println("")

	switch choice {
	case "y":
		fmt.Println("Continuing")
		fmt.Println("")
	case "n":
		fmt.Println("Exiting..")
		os.Exit(0)
	default:
		fmt.Println("Exiting. Unknown input")
		os.Exit(1)
	}

	var cycleCount int
	fmt.Println("how many times you want to perform this operation (and check fluctuations in the weighted mean fluctuation)")
	fmt.Scanf("%d", &cycleCount)
	fmt.Println("")

	for i := 0; i < cycleCount; i++ {

		var addOrDel string
		fmt.Println("Do you want to add/update an customer or delete an exising one ?")
		fmt.Println("Enter '1' to add or update")
		fmt.Println("Enter '2' to delete existing customer")
		fmt.Scanf("%s", &addOrDel)
		fmt.Println("")

		switch addOrDel {
		case "1":
			var customerName string
			fmt.Println("Enter customer name:")
			fmt.Scanf("%s", &customerName)
			fmt.Println("")

			var newCoreCount int
			fmt.Println("Enter core count:")
			fmt.Scanf("%d", &newCoreCount)
			fmt.Println("")

			if _, ok := dailyCoreConsumptionPerCustomer[customerName]; !ok {
				fmt.Printf("New customer %s. Adding to the customer List", customerName)
				fmt.Println("")
				customersList = append(customersList, customerName)
			}
			dailyCoreConsumptionPerCustomer[customerName] = newCoreCount
			generateWeightedMean()
			fmt.Println("")

		case "2":
			var customerName string
			fmt.Println("Enter customer name:")
			fmt.Scanf("%s", &customerName)
			fmt.Println("")

			if _, ok := dailyCoreConsumptionPerCustomer[customerName]; !ok {
				fmt.Printf("Customer %s does not exist. Wrong input", customerName)
				fmt.Println("")
			} else {
				customersList = deleteValueFromArray(customersList, customerName)
				delete(dailyCoreConsumptionPerCustomer, customerName)
				generateWeightedMean()
				fmt.Println("")
			}
		default:
			fmt.Println("invalid input")
			os.Exit(1)
		}
	}
	createLineChart()
	fmt.Println("A line chart has been created in the current directory")
}

func generateWeightedMean() {
	// Calculate the sum of the daily core consumption per customer
	sumOfDailyCoreConsumptionPerCustomer := 0
	maxValue := -1
	for _, value := range dailyCoreConsumptionPerCustomer {
		sumOfDailyCoreConsumptionPerCustomer += value
		if value > maxValue {
			maxValue = value
		}
	}
	fmt.Println("Total number of Cores Consumed:", sumOfDailyCoreConsumptionPerCustomer)
	fmt.Println("")

	// Calculate the weights for each customer
	weightsSum := 0.0
	weights := make(map[string]float64, len(customersList))
	for _, customerName := range customersList {
		//weights[customerName] = float64(dailyCoreConsumptionPerCustomer[customerName]) / float64(sumOfDailyCoreConsumptionPerCustomer)
		//weights[customerName] = float64(dailyCoreConsumptionPerCustomer[customerName]) / float64(maxValue)
		weights[customerName] = float64(dailyCoreConsumptionPerCustomer[customerName]) / float64(maxNumberOfCoresThatCanBeConsumed)
		weightsSum += weights[customerName]
	}
	fmt.Println("Daily Core Consumption Per Customer and their assigned weights")
	displayMap(dailyCoreConsumptionPerCustomer, weights, customersList)
	fmt.Println("")

	fmt.Println("sum of the weights: ", weightsSum)
	fmt.Println("")

	// Recalculate the weighted average mean of total cores consumed with weights from 0 to 1
	recalculatedWeightedAverageMean := 0.0
	for _, customerName := range customersList {
		recalculatedWeightedAverageMean += float64(dailyCoreConsumptionPerCustomer[customerName]) * weights[customerName]
	}

	//recalculatedWeightedAverageMean = float64(recalculatedWeightedAverageMean) / float64(sumOfDailyCoreConsumptionPerCustomer) * 100
	// Print the recalculated weighted average mean
	fmt.Println("Recalculated weighted average mean of total cores consumed with weights from 0 to 1:", recalculatedWeightedAverageMean)
	fmt.Println("")

	newCoordinate := Coordinate{
		TotalCores:      sumOfDailyCoreConsumptionPerCustomer,
		WeightedAverage: recalculatedWeightedAverageMean,
		TotalCustomers:  len(customersList),
		ActionCount:     globalCounter,
	}
	myPlot = append(myPlot, newCoordinate)
	globalCounter++
}

func deleteValueFromArray(array []string, value string) []string {
	// Find the index of the value you want to delete.
	index := -1
	for i, v := range array {
		if v == value {
			index = i
			break
		}
	}

	// If the value is not in the array, return the original array.
	if index == -1 {
		return array
	}

	// Create a new array with a capacity of one less than the original array.
	newArray := make([]string, len(array)-1)

	// Iterate over the original array and copy all of the elements to the new array, except for the element at the index of the value you want to delete.
	for i := 0; i < index; i++ {
		newArray[i] = array[i]
	}
	for i := index + 1; i < len(array); i++ {
		newArray[i-1] = array[i]
	}

	// Return the new array.
	return newArray
}

func displayMap[K comparable, V any, U any](customerCores map[K]V, weights map[K]U, customersList []K) {
	for _, customerName := range customersList {
		fmt.Printf("%s: %v: %v\n", customerName, customerCores[customerName], weights[customerName])
	}
}

func makeCustomerSortedList(numberOfCustomers int) []string {
	// Extract keys from the map
	var keys []string
	for i := 1; i <= numberOfCustomers; i++ {
		if i < 10 {
			keys = append(keys, "Customer_0"+strconv.Itoa(i))
		} else {
			keys = append(keys, "Customer_"+strconv.Itoa(i))
		}
	}

	// Sort the keys
	sort.Strings(keys)

	return keys
}

func generateWeightedAverage() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(myPlot); i++ {
		items = append(items, opts.LineData{Value: myPlot[i].WeightedAverage, Name: "WeightedAverage"})
	}
	return items
}
func generateTotalCores() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(myPlot); i++ {
		items = append(items, opts.LineData{Value: myPlot[i].TotalCores, Name: "TotalCores"})
	}
	return items
}
func generateTotalCustomers() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(myPlot); i++ {
		items = append(items, opts.LineData{Value: myPlot[i].TotalCustomers, Name: "TotalCustomers"})
	}
	return items
}

func createLineChart() {
	// create a new line instance
	line := charts.NewLine()

	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeInfographic,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line chart representing weighted moving average",
			Subtitle: "Notice the fluctuations in the weighted mean as the number of cores changes",
		}),
	)

	// Put data into instance
	line.SetXAxis([]string{"Oct", "Nov", "Dec", "Jan", "Feb", "Mar"}).
		AddSeries("Category A", generateWeightedAverage()).
		AddSeries("Category B", generateTotalCores()).
		AddSeries("Category C", generateTotalCustomers()).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create("line.html")
	_ = line.Render(f)
}
