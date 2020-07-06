package example

default can_earn_cashback = false
default cashback_amount_percent = 2.5

# This should be in data...so that it won't appear in the result.
# {
#     "rates": [
# 		{ "country": "*", "rate": 10 }
# 	]
# }

# Input:
# {
#     "last_purchased_at": "2020-07-06",
#     "sku": "product-123",
#     "country": "singapore"
# }


# Function to count elapsed days.
elapsed_days (yyyymmdd) = days {
	date := time.parse_ns("2006-01-02", yyyymmdd)
    now := time.now_ns()
    elapsed_ns := now - date
    elapsed_s := elapsed_ns / 1e9
    days := elapsed_s / 86400
}

can_earn_cashback {
    sku := input.sku
    sku == "product-123"
}

default new_purchase_last_ten_days = false

new_purchase_last_ten_days {
    days := elapsed_days(input.last_purchased_at)
	days < 10
}

cashback_amount_percent = 5 {
	can_earn_cashback
    not new_purchase_last_ten_days
}

cashback_amount_percent = 10 {
	can_earn_cashback
    new_purchase_last_ten_days
    input.country == "malaysia"
}


cashback_amount_percent = 7.5 {
	can_earn_cashback
    new_purchase_last_ten_days
    input.country == "singapore"
}

cashback_amount_percent = z {
	can_earn_cashback
    new_purchase_last_ten_days
	row := data.rates[_]
    row.country == input.country
    z := row.rate * 2
}
