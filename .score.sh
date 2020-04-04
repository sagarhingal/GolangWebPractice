# Build and run solution
rm -f question
go mod vendor
go build -o question
chmod +x question
./question &>/dev/null &
SERVER_PID=$!
# Setup test environment
cd .questiontest
rm -f -r ./venv
virtualenv -p /usr/bin/python3 venv
source venv/bin/activate
pip install -r requirements.txt
sleep 5
py.test --suppress-tests-failed-exit-code --junitxml=results.xml
kill $SERVER_PID
sleep 1
python score.py
