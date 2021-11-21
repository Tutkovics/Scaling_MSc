$ echo "GET http://152.66.211.2:30000/forward?" | vegeta attack -duration=600s -rate=70 -timeout=5s -max-workers=1000 -workers=1000 | tee results.bin | vegeta rep
ort -type=json  > ../results/test-HPA-1110/vegeta-results.json

