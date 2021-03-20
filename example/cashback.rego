package example

# Test it here.
# https://play.openpolicyagent.org/

# This should be in data...so that it won't appear in the result.
#{
    #"rates": [
     #{ "country": "*", "rate": 10 }
     ##{ "country": "singapore", "rate": 5 } policy.rego:40: eval_conflict_error: complete rules must not produce multiple outputs
   #],
    #"products": [
        #{ "sku": "product-1", "cashback_elligible": true },
        #{ "sku": "product-2", "cashback_elligible": false }
    #]
#}

# Input:
#{
    #"last_purchased_at": "2021-03-20",
    #"sku": "product-1",
    #"country": "singapore"
#}

#Output:
#{
    #"can_earn_cashback": true,
    #"cashback_amount_percent": 7.5,
    #"purchase_within_last_ten_days": true
#}


# Function to count elapsed days.
elapsed_days (yyyymmdd) = days {
	date := time.parse_ns("2006-01-02", yyyymmdd)
    now := time.now_ns()
    elapsed_ns := now - date
    elapsed_s := elapsed_ns / 1e9
    days := elapsed_s / 86400
}

# Provide a default value, so that if it fails to evaluate, there will always
# still be a result.
default can_earn_cashback = false

# READ: If the product sku matches that from data, and is cashback_elligible,
# returns true.
can_earn_cashback {
    product := data.products[_]
    product.sku = input.sku
    product.cashback_elligible
}

default purchase_within_last_ten_days = false

purchase_within_last_ten_days {
    days := elapsed_days(input.last_purchased_at)
    days < 10
}

# If none matches, the value will be 2.5%.
default cashback_amount_percent = 2.5

# Each of this is an OR statement. Whichever matches first will return.
# READ: If the product is cashback elligible, and is not purchase within the
# last ten days, returns 5 percent.
cashback_amount_percent = 5 {
    # Each line within is an AND statement.
    can_earn_cashback
    not purchase_within_last_ten_days
}

# READ: If the product is cashback elligible, and is purchased within the last
# ten days, and the country is Malaysia, returns 10 percent.
cashback_amount_percent = 10 {
    can_earn_cashback
    purchase_within_last_ten_days
    input.country == "malaysia"
}

# READ: If the product is cashback elligible, and is purchased within the last
# ten days, and the country is Singapore, returns 7.5 percent.
cashback_amount_percent = 7.5 {
    can_earn_cashback
    purchase_within_last_ten_days
    input.country == "singapore"
}

# READ: If the product is cashback elligible, and is purchased within the last
# ten days, and the country is in the list of rates, returns z percent, where z
# is the rate for that given country.
cashback_amount_percent = z {
    can_earn_cashback
    purchase_within_last_ten_days
    rate := data.rates[_]
    rate.country == input.country
    z := rate.rate * 2
}
