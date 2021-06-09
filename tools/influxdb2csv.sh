#!/usr/bin/env bash

function print_usage_help_and_exit
{
    echo $(basename $0) "-u USER -p PASSWORD -f FILENAME"
    cat <<_EOF
  -u USER: influxdb username
  -p PASSWORD: influxdb password
  -c FILENAME: filename to be used for storing data
  -h: print usage

ENV:
    INFLUX_USER - default for the parameter -u
    INFLUX_PASSWORD - default for the parameter -p
    DATA_FILENAME - default for the parameter -f
_EOF
    exit $1
}

while getopts ":u:p:f:h" OPTION; do
    case $OPTION in
    u)
        INFLUX_USER=$OPTARG
        ;;
    p)
        INFLUX_PASSWORD=$OPTARG
        ;;
    f)
        DATA_FILENAME=$OPTARG
        ;;
    h)
        print_usage_help_and_exit 0
        ;;
    *)
        echo "Incorrect option provided"
        print_usage_help_and_exit 1
        exit 1
        ;;
    esac
done

DATA_FILENAME=${DATA_FILENAME:-data}

if [ -z "${INFLUX_USER:-}" ]; then
    echo "InfluxDB username not specified!"
    echo "Either use the parameter -u or set INFLUX_USER env variable"
    echo "Run the script with -h for more help"
    exit 1
fi

if [ -z "${INFLUX_PASSWORD:-}" ]; then
    echo "InfluxDB password not specified!"
    echo "Either use the parameter -p or set INFLUX_PASSWORD env variable"
    echo "Run the script with -h for more help"
    exit 1
fi

URL='https://iot2.pavoucek.net/influxdb/query?pretty=true'

TIMEFROM='2021-05-31T00:00:00Z'
TIMETO='2021-06-07T00:00:00Z'

ID_B3012_TEMP='60a17fb4b3c965886d4c49f1'
ID_B3013_TEMP='60a2772bb3c965886d4c49f4'

ID_B3012_HUM='60a18065b3c965886d4c49f2'
ID_B3013_HUM='60a2770eb3c965886d4c49f3'

ID_B3014_TEMP='60a2ce0ab3c965886d4c49f5'
ID_B3015_TEMP='B3015_Temp'
ID_B3017_TEMP='B3017_Temp'
ID_B3018_TEMP='60a509c7b3c965886d4c4a05'
ID_B3019_TEMP='60a509cfb3c965886d4c4a06'

ID_B3014_HUM='60a2ce23b3c965886d4c49f6'
ID_B3017_HUM='60a36f24b3c965886d4c4a01'
ID_B3018_HUM='60a509e7b3c965886d4c4a07'
ID_B3019_HUM='60a509f1b3c965886d4c4a08'

QUERY="q=SELECT MEAN(\"value\") FROM \"sensor\" WHERE time >= '$TIMEFROM' AND time <= '$TIMETO' AND \"id\" = '$ID_B3019_HUM' GROUP BY time(1h)"

#QUERY="q=SELECT MEAN(\"value\") FROM \"sensor\" WHERE time >= '$TIMEFROM' AND time <= '$TIMETO' AND (\"id\" = '$ID_B3012_TEMP' OR \"id\" = '$ID_B3013_TEMP') GROUP BY time(1h), \"id\""

echo "Fetching InfluxDB data to file $DATA_FILENAME.json"
# tip:  -H "Accept: application/csv"
curl -G -u $INFLUX_USER:$INFLUX_PASSWORD \
    $URL \
    -o $DATA_FILENAME.json \
    --data-urlencode "db=surgal" \
    --data-urlencode "precision=m" \
    --data-urlencode "$QUERY"

echo "Transforming $DATA_FILENAME.json -> $DATA_FILENAME.json"
cat $DATA_FILENAME.json | jq -r ".results[0].series[0].values" > $DATA_FILENAME.trans.json

echo "Cnverting $DATA_FILENAME.json -> $DATA_FILENAME.csv"
cat $DATA_FILENAME.json | jq -r "(.results[0].series[0].columns), (.results[0].series[0].values[]) | @csv" > $DATA_FILENAME.csv

echo "Done"
