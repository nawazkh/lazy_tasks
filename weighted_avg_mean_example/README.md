# Weighted Average Mean Example

This code demonstrates that weighted Average mean is good metric to track an entity in case of an imbalanced dataset.

- This demo will output a `line.html` which can be opened in a browser to see the results.

## Demo run

```bash
$ go run weighted_avg_mean.go

Enter initial number of customers to consider
20

Enter max number of cores that can be consumed by a customer
1000

Total number of Cores Consumed: 11668

Daily Core Consumption Per Customer and their assigned weights
Customer_01: 455: 0.455
Customer_02: 849: 0.849
Customer_03: 68: 0.068
Customer_04: 940: 0.94
Customer_05: 429: 0.429
Customer_06: 345: 0.345
Customer_07: 540: 0.54
Customer_08: 56: 0.056
Customer_09: 536: 0.536
Customer_10: 619: 0.619
Customer_11: 922: 0.922
Customer_12: 734: 0.734
Customer_13: 327: 0.327
Customer_14: 912: 0.912
Customer_15: 885: 0.885
Customer_16: 95: 0.095
Customer_17: 863: 0.863
Customer_18: 854: 0.854
Customer_19: 954: 0.954
Customer_20: 285: 0.285

sum of the weights:  11.668

Recalculated weighted average mean of total cores consumed with weights from 0 to 1: 8669.498

Would you like to add/delete customers? y/n
y

Continuing

how many times you want to perform this operation (and check fluctuations in the weighted mean fluctuation)
5

Do you want to add/update an customer or delete an exising one ?
Enter '1' to add or update
Enter '2' to delete existing customer
1

Enter customer name:
Customer_11

Enter core count:
50

Total number of Cores Consumed: 10796

Daily Core Consumption Per Customer and their assigned weights
Customer_01: 455: 0.455
Customer_02: 849: 0.849
Customer_03: 68: 0.068
Customer_04: 940: 0.94
Customer_05: 429: 0.429
Customer_06: 345: 0.345
Customer_07: 540: 0.54
Customer_08: 56: 0.056
Customer_09: 536: 0.536
Customer_10: 619: 0.619
Customer_11: 50: 0.05
Customer_12: 734: 0.734
Customer_13: 327: 0.327
Customer_14: 912: 0.912
Customer_15: 885: 0.885
Customer_16: 95: 0.095
Customer_17: 863: 0.863
Customer_18: 854: 0.854
Customer_19: 954: 0.954
Customer_20: 285: 0.285

sum of the weights:  10.796

Recalculated weighted average mean of total cores consumed with weights from 0 to 1: 7821.914


Do you want to add/update an customer or delete an exising one ?
Enter '1' to add or update
Enter '2' to delete existing customer
1

Enter customer name:
Customer_14

Enter core count:
10

Total number of Cores Consumed: 9894

Daily Core Consumption Per Customer and their assigned weights
Customer_01: 455: 0.455
Customer_02: 849: 0.849
Customer_03: 68: 0.068
Customer_04: 940: 0.94
Customer_05: 429: 0.429
Customer_06: 345: 0.345
Customer_07: 540: 0.54
Customer_08: 56: 0.056
Customer_09: 536: 0.536
Customer_10: 619: 0.619
Customer_11: 50: 0.05
Customer_12: 734: 0.734
Customer_13: 327: 0.327
Customer_14: 10: 0.01
Customer_15: 885: 0.885
Customer_16: 95: 0.095
Customer_17: 863: 0.863
Customer_18: 854: 0.854
Customer_19: 954: 0.954
Customer_20: 285: 0.285

sum of the weights:  9.893999999999998

Recalculated weighted average mean of total cores consumed with weights from 0 to 1: 6990.2699999999995


Do you want to add/update an customer or delete an exising one ?
Enter '1' to add or update
Enter '2' to delete existing customer
1

Enter customer name:
Test_user

Enter core count:
500

New customer Test_user. Adding to the customer List
Total number of Cores Consumed: 10394

Daily Core Consumption Per Customer and their assigned weights
Customer_01: 455: 0.455
Customer_02: 849: 0.849
Customer_03: 68: 0.068
Customer_04: 940: 0.94
Customer_05: 429: 0.429
Customer_06: 345: 0.345
Customer_07: 540: 0.54
Customer_08: 56: 0.056
Customer_09: 536: 0.536
Customer_10: 619: 0.619
Customer_11: 50: 0.05
Customer_12: 734: 0.734
Customer_13: 327: 0.327
Customer_14: 10: 0.01
Customer_15: 885: 0.885
Customer_16: 95: 0.095
Customer_17: 863: 0.863
Customer_18: 854: 0.854
Customer_19: 954: 0.954
Customer_20: 285: 0.285
Test_user: 500: 0.5

sum of the weights:  10.393999999999998

Recalculated weighted average mean of total cores consumed with weights from 0 to 1: 7240.2699999999995


Do you want to add/update an customer or delete an exising one ?
Enter '1' to add or update
Enter '2' to delete existing customer
2

Enter customer name:
Customer_19

Total number of Cores Consumed: 9440

Daily Core Consumption Per Customer and their assigned weights
Customer_01: 455: 0.455
Customer_02: 849: 0.849
Customer_03: 68: 0.068
Customer_04: 940: 0.94
Customer_05: 429: 0.429
Customer_06: 345: 0.345
Customer_07: 540: 0.54
Customer_08: 56: 0.056
Customer_09: 536: 0.536
Customer_10: 619: 0.619
Customer_11: 50: 0.05
Customer_12: 734: 0.734
Customer_13: 327: 0.327
Customer_14: 10: 0.01
Customer_15: 885: 0.885
Customer_16: 95: 0.095
Customer_17: 863: 0.863
Customer_18: 854: 0.854
Customer_20: 285: 0.285
Test_user: 500: 0.5

sum of the weights:  9.439999999999998

Recalculated weighted average mean of total cores consumed with weights from 0 to 1: 6330.1539999999995


Do you want to add/update an customer or delete an exising one ?
Enter '1' to add or update
Enter '2' to delete existing customer
2

Enter customer name:
Customer_07

Total number of Cores Consumed: 8900

Daily Core Consumption Per Customer and their assigned weights
Customer_01: 455: 0.455
Customer_02: 849: 0.849
Customer_03: 68: 0.068
Customer_04: 940: 0.94
Customer_05: 429: 0.429
Customer_06: 345: 0.345
Customer_08: 56: 0.056
Customer_09: 536: 0.536
Customer_10: 619: 0.619
Customer_11: 50: 0.05
Customer_12: 734: 0.734
Customer_13: 327: 0.327
Customer_14: 10: 0.01
Customer_15: 885: 0.885
Customer_16: 95: 0.095
Customer_17: 863: 0.863
Customer_18: 854: 0.854
Customer_20: 285: 0.285
Test_user: 500: 0.5

sum of the weights:  8.899999999999999

Recalculated weighted average mean of total cores consumed with weights from 0 to 1: 6038.554


A line chart has been created in the current directory
```

## Output line chart

![line chart](./weighted_avg_image.png)
