google.charts.load('current', { 'packages': ['corechart'] });
google.charts.setOnLoadCallback(drawChart);

var priceData;
var emaData;
var macdData;
var priceChart;
var emaChart;
var macdChart;
var socket;

function drawChart() {
  socket   = new WebSocket('ws://localhost:3000/price');
  socket.onopen = function (event) {
    console.log('Connected to server');
  };

  socket.onmessage = function (event) {
    var chartData = JSON.parse(event.data);

    // Convert price and ema to float numbers
    var price = parseFloat(chartData.price_data);

    console.log("Raw Tech data:", chartData.technical); 

    // Technical data
    var shortEMA = parseFloat(chartData.technical.short_ema);
    var longEMA = parseFloat(chartData.technical.long_ema);
    var macdLine = parseFloat(chartData.technical.macd_line);
    var signalLine = parseFloat(chartData.technical.signal_ema);
    var histogram = parseFloat(chartData.technical.histogram);

    // ChartDataから価格、EMAとタイムスタンプ（またはTick IDなど）を取得し、グラフに追加
    addPriceData(chartData.timestamp, price);
    addEMAData(chartData.timestamp, shortEMA,longEMA);
    addMACDData(chartData.timestamp, macdLine,signalLine,histogram);
  };

  priceChart = new google.visualization.LineChart(document.getElementById('price_chart'));
  emaChart = new google.visualization.LineChart(document.getElementById('ema_chart'));
  macdChart = new google.visualization.LineChart(document.getElementById('macd_chart'));

  socket.onclose = function (event) {
    console.log('Disconnected from server');
  };

  socket.onerror = function (error) {
    console.log('Error occurred:', error);
  };
  
  priceData = google.visualization.arrayToDataTable([
    ["date","price"], 
    ['0',4000000] 
  ]);

  emaData = google.visualization.arrayToDataTable([
    ["date","shortEMA","longEMA"],
    ['0',0,0]
  ]);

  macdData = google.visualization.arrayToDataTable([
    ["date","MACD","Signal","Histogram"],
    ['0',0,0,0]
  ])



  var options = {
    title: 'Price Data',
    curveType: 'function',
    legend: { position: 'bottom' }
  };

  priceChart.draw(priceData, options);
  priceData.removeRow(0);

  emaChart.draw(emaData, options);
  emaData.removeRow(0);
  macdChart.draw(macdData, options);
  macdData.removeRow(0);
}

function addPriceData(tick, price) {
  var dateTick = new Date(tick);
  if (isNaN(dateTick)) { // 日付が無効な場合は処理を中断
    console.log('無効な日付です:', tick);
    return;
  }
  var formattedDate = ("00"+dateTick.getHours()).slice(-2)+":"+("00"+dateTick.getMinutes()).slice(-2)+":"+("00"+dateTick.getSeconds()).slice(-2); 
  if (priceData.getNumberOfRows() > 200) {
    priceData.removeRow(0);
  }
  console.log(formattedDate + ":" + price);
  priceData.addRow([formattedDate, price]); 
  priceChart.draw(priceData, null);
}

function addEMAData(tick, short,long) {
  var dateTick = new Date(tick);
  if (isNaN(dateTick)) { // 日付が無効な場合は処理を中断
    console.log('無効な日付です:', tick);
    return;
  }
  var formattedDate = ("00"+dateTick.getHours()).slice(-2)+":"+("00"+dateTick.getMinutes()).slice(-2)+":"+("00"+dateTick.getSeconds()).slice(-2); 
  if (emaData.getNumberOfRows() > 200) {
    emaData.removeRow(0);
  }
  console.log(formattedDate + ":" + short + ":" + long);
  emaData.addRow([formattedDate, short,long]); 
  emaChart.draw(emaData, null); 
}

function addMACDData(tick, macd,signal,histogram) {
    var dateTick = new Date(tick);
    if (isNaN(dateTick)) { // 日付が無効な場合は処理を中断
      console.log('無効な日付です:', tick);
      return;
    }
    var formattedDate = ("00"+dateTick.getHours()).slice(-2)+":"+("00"+dateTick.getMinutes()).slice(-2)+":"+("00"+dateTick.getSeconds()).slice(-2); 
    if (macdData.getNumberOfRows() > 200) {
      macdData.removeRow(0);
    }
    console.log(formattedDate + ":" + macd + ":" + signal + ":" + histogram);
    macdData.addRow([formattedDate, macd,signal,histogram]); 
    macdChart.draw(macdData, null); 
}

window.addEventListener('beforeunload', function (event) {
  socket.close();
});

