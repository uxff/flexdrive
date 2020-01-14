<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0"/>
    <title>{{.appname}}</title>
    <link href="/static/bower/bootstrap/dist/css/bootstrap.min.css" rel="stylesheet">
    {{range .HeadStyles}}
        <link rel="stylesheet" href="{{.}}">
    {{end}}
</head>
<body id="home">

  {{.BaseHeader}}
  
  <div id="wrap">
      {{.LayoutContent}}
  </div>
  
  {{.BaseFooter}}

  <script src="/static/bower/jquery/dist/jquery.min.js"></script>
  <script src="/static/bower/bootstrap/dist/js/bootstrap.min.js" type="text/javascript" charset="utf-8"></script>
  <script type="text/javascript" src="/static/js/twitter-bootstrap-hover-dropdown.min.js"></script>
  {{range .HeadScripts}}
      <script src="{{.}}"></script>
  {{end}}
 
</body>
</html>
