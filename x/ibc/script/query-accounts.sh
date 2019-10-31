gaiacli --home ibc-gaia/n0/gaiacli q account -o text $(gaiacli --home ibc-gaia/n0/gaiacli keys show n0 | jq -r '.address')
iriscli --home ibc-iris/n0/iriscli q account -o text $(iriscli --home ibc-iris/n0/iriscli keys show n0 | jq -r '.address')
