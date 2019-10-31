sleep 3
echo "" && echo "### create-client ..." && echo ""
gaiacli --home ibc-gaia/n0/gaiacli q ibc client self-consensus-state -o json >ibc-iris/n0/consensus_state.json
echo 12345678 | iriscli --home ibc-iris/n0/iriscli tx ibc client create clientiris ibc-iris/n0/consensus_state.json --from n0 -y --broadcast-mode=block
iriscli --home ibc-iris/n0/iriscli q ibc client self-consensus-state -o json >ibc-gaia/n0/consensus_state.json
echo 12345678 | gaiacli --home ibc-gaia/n0/gaiacli tx ibc client create clientgaia ibc-gaia/n0/consensus_state.json --from n0 -y --broadcast-mode=block
echo "" && echo "### query-client ..." && echo ""
iriscli --home ibc-iris/n0/iriscli q ibc client state clientiris | jq
gaiacli --home ibc-gaia/n0/gaiacli q ibc client state clientgaia | jq

echo "" && echo "### open-init ..." && echo ""
echo 12345678 | iriscli --home ibc-iris/n0/iriscli tx ibc connection open-init connectionid clientiris connectionid clientgaia prefix.json --from n0 -y --broadcast-mode=block
echo "" && echo "### open-try ..." && echo ""
sleep 3 && iriscli --home ibc-iris/n0/iriscli q ibc client header -o json >ibc-gaia/n0/header.json
iriscli --home ibc-iris/n0/iriscli q ibc connection proof connectionid $(jq -r '.value.SignedHeader.header.height' ibc-gaia/n0/header.json) -o json >ibc-gaia/n0/conn_proof_init.json
echo 12345678 | gaiacli --home ibc-gaia/n0/gaiacli tx ibc client update clientgaia ibc-gaia/n0/header.json --from n0 -y --broadcast-mode=block
echo 12345678 | gaiacli --home ibc-gaia/n0/gaiacli tx ibc connection open-try connectionid clientgaia connectionid clientiris prefix.json 1.0.0 ibc-gaia/n0/conn_proof_init.json $(jq -r '.value.SignedHeader.header.height' ibc-gaia/n0/header.json) --from n0 -y --broadcast-mode=block
echo "" && echo "### open-ack ..." && echo ""
sleep 3 && gaiacli --home ibc-gaia/n0/gaiacli q ibc client header -o json >ibc-iris/n0/header.json
gaiacli --home ibc-gaia/n0/gaiacli q ibc connection proof connectionid $(jq -r '.value.SignedHeader.header.height' ibc-iris/n0/header.json) -o json >ibc-iris/n0/conn_proof_try.json
echo 12345678 | iriscli --home ibc-iris/n0/iriscli tx ibc client update clientiris ibc-iris/n0/header.json --from n0 -y --broadcast-mode=block
echo 12345678 | iriscli --home ibc-iris/n0/iriscli tx ibc connection open-ack connectionid ibc-iris/n0/conn_proof_try.json $(jq -r '.value.SignedHeader.header.height' ibc-iris/n0/header.json) 1.0.0 --from n0 -y --broadcast-mode=block
echo "" && echo "### open-confirm ..." && echo ""
sleep 3 && iriscli --home ibc-iris/n0/iriscli q ibc client header -o json >ibc-gaia/n0/header.json
iriscli --home ibc-iris/n0/iriscli q ibc connection proof connectionid $(jq -r '.value.SignedHeader.header.height' ibc-gaia/n0/header.json) -o json >ibc-gaia/n0/conn_proof_ack.json
echo 12345678 | gaiacli --home ibc-gaia/n0/gaiacli tx ibc client update clientgaia ibc-gaia/n0/header.json --from n0 -y --broadcast-mode=block
echo 12345678 | gaiacli --home ibc-gaia/n0/gaiacli tx ibc connection open-confirm connectionid ibc-gaia/n0/conn_proof_ack.json $(jq -r '.value.SignedHeader.header.height' ibc-gaia/n0/header.json) --from n0 -y --broadcast-mode=block
echo "" && echo "### query-connection ..." && echo ""
iriscli --home ibc-iris/n0/iriscli q ibc connection end connectionid | jq
gaiacli --home ibc-gaia/n0/gaiacli q ibc connection end connectionid | jq

echo "" && echo "### open-init ..." && echo ""
echo 12345678 | iriscli --home ibc-iris/n0/iriscli tx ibc channel open-init ppppppport channeliris ppppppport channelgaia connectionid --ordered=false --from n0 -y --broadcast-mode=block
echo "" && echo "### open-try ..." && echo ""
sleep 3 && iriscli --home ibc-iris/n0/iriscli q ibc client header -o json >ibc-gaia/n0/header.json
iriscli --home ibc-iris/n0/iriscli q ibc channel proof ppppppport channeliris $(jq -r '.value.SignedHeader.header.height' ibc-gaia/n0/header.json) -o json >ibc-gaia/n0/chann_proof_init.json
echo 12345678 | gaiacli --home ibc-gaia/n0/gaiacli tx ibc client update clientgaia ibc-gaia/n0/header.json --from n0 -y --broadcast-mode=block
echo 12345678 | gaiacli --home ibc-gaia/n0/gaiacli tx ibc channel open-try ppppppport channelgaia ppppppport channeliris connectionid ibc-gaia/n0/chann_proof_init.json $(jq -r '.value.SignedHeader.header.height' ibc-gaia/n0/header.json) --ordered=false --from n0 -y --broadcast-mode=block
echo "" && echo "### open-ack ..." && echo ""
sleep 3 && gaiacli --home ibc-gaia/n0/gaiacli q ibc client header -o json >ibc-iris/n0/header.json
gaiacli --home ibc-gaia/n0/gaiacli q ibc channel proof ppppppport channelgaia $(jq -r '.value.SignedHeader.header.height' ibc-iris/n0/header.json) -o json >ibc-iris/n0/chann_proof_try.json
echo 12345678 | iriscli --home ibc-iris/n0/iriscli tx ibc client update clientiris ibc-iris/n0/header.json --from n0 -y --broadcast-mode=block
echo 12345678 | iriscli --home ibc-iris/n0/iriscli tx ibc channel open-ack ppppppport channeliris ibc-iris/n0/chann_proof_try.json $(jq -r '.value.SignedHeader.header.height' ibc-iris/n0/header.json) --from n0 -y --broadcast-mode=block
echo "" && echo "### open-confirm ..." && echo ""
sleep 3 && iriscli --home ibc-iris/n0/iriscli q ibc client header -o json >ibc-gaia/n0/header.json
iriscli --home ibc-iris/n0/iriscli q ibc channel proof ppppppport channeliris $(jq -r '.value.SignedHeader.header.height' ibc-gaia/n0/header.json) -o json >ibc-gaia/n0/chann_proof_ack.json
echo 12345678 | gaiacli --home ibc-gaia/n0/gaiacli tx ibc client update clientgaia ibc-gaia/n0/header.json --from n0 -y --broadcast-mode=block
echo 12345678 | gaiacli --home ibc-gaia/n0/gaiacli tx ibc channel open-confirm ppppppport channelgaia ibc-gaia/n0/chann_proof_ack.json $(jq -r '.value.SignedHeader.header.height' ibc-gaia/n0/header.json) --from n0 -y --broadcast-mode=block
echo "" && echo "### query-channel ..." && echo ""
iriscli --home ibc-iris/n0/iriscli query ibc channel end ppppppport channeliris | jq
gaiacli --home ibc-gaia/n0/gaiacli query ibc channel end ppppppport channelgaia | jq
