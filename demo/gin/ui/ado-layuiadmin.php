<?php
$content = file_get_contents('index.html');
//<a lay-href="home/console.html">控制台</a>

$pattern='/lay-href="([^"]*)"/i';
preg_match_all($pattern,$content,$arr);

$arr[1]=[
    'app/content/listform.html',
    'app/forum/replysform.html',
    'app/forum/listform.html',
    'app/content/tagsform.html',
    'app/content/contform.html',
    'app/workorder/listform.html',
    'user/user/userform.html',
    'user/administrators/adminform.html',
    'user/administrators/roleform.html',

];

$urls = [];
$url_prefix='https://www.layui.com/admin/std/dist/views/';
foreach($arr[1] as $k=>$url){
    if(strstr($url,'www')){
        continue;
    }
    $tmp = strripos($url,'/');
    $file_path = substr($url,0,$tmp);
    $file_name = substr($url,$tmp);
    if(!is_dir($file_path)){
        mkdir($file_path,0777,true);
    }
    file_put_contents($file_path.$file_name,file_get_contents($url_prefix.$url));
}
