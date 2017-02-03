logid=1000
var router = new Navigo();

router.on('/buckets', function () {
    loadBucketTable();
    $('#pg1').hide();
    $('#pg3').hide();
    $('#pg4').hide();
    $('#pg2').show();
});

router.on('/prefixScan', function () {
    $('#pg1').hide();
    $('#pg2').hide();
    $('#pg4').hide();
    $('#pg3').show();
});

router.on('/explore', function () {
    $('#pg1').hide();
    $('#pg2').hide();
    $('#pg3').hide();
    $('#pg4').show();
});

router.on('/', function () {
    $('#pg2').hide();
    $('#pg3').hide();
    $('#pg4').hide();
    $('#pg1').show();
});

router.on(function() {
    $('#pg2').hide();
    $('#pg3').hide();
    $('#pg4').hide();

    console.log("default route:no other routes matched.")
});

function doDelete(key){
    var r = confirm("Delete?");

    if (r == true) {
        b = $('#pbucket').val();
        deleteKeyReq(b,key);
        window.setTimeout(prefixScan, 1000);
    }
}

function doEdit(key){
    b = $('#pbucket').val();
    getRequest(b,key);

    $('#bucket').val(b);
    $('#key').val(key);
    router.navigate('#/');
}

function doPrefixScan(bucket){
    $('#pbucket').val(bucket);
    $('#pkey').val("");

    prefixScan()

    router.navigate('#/prefixScan');
}

function doExplore(bucket) {
    $('#ebucket').val(bucket);

    explore();

    router.navigate('#/explore');
}

function log(text){
    console.log(text)
    $('#log').append("<br/>["+logid+"] "+ JSON.stringify(text))

    logid++
    $('#log').scrollTop($('#log')[0].scrollHeight - $('#log')[0].clientHeight);
}

function get(){
    getRequest($('#bucket').val(),$('#key').val());
}

function getRequest(bucket,key){
    $.post("/get",{bucket:bucket,key:key},function(data){
        log(data)
        if(data[0]=="ok"){
            $('#value').val(data[1])
        }
    });
}

function deleteBucket(){
    $.post("/deleteBucket",{bucket:$('#bucket').val()},function(data){
        log(data)
    });
}

function deleteKey(){
    deleteKeyReq($('#bucket').val(),$('#key').val());
}

function deleteKeyReq(bucket,key){
    $.post("/deleteKey",{bucket:bucket,key:key},function(data){
        log(data)
    });

}

function put(){
    $.post("/put",{bucket:$('#bucket').val(),key:$('#key').val(),value:$('#value').val()},function(data){
        log(data)
    });
}

function prefixScan() {
    $('#pfs').html("")
    var source = $('#prefixtpl').html();
    var template = Handlebars.compile(source);

    $.post("/prefixScan",{bucket:$('#pbucket').val(),key:$('#pkey').val()},function(data){
        log(data)
        //var rendered = Mustache.render(template, {list: data.M});
        var html    = template({list: data.M});
        $('#pfs').html(html)
    });
}

function explore() {
    $('#expl').html("");
    $('#bdc').html("");

     var source = $('#exploretpl').html();
     var template = Handlebars.compile(source);

     source = $('#breadcrumbtpl').html();
     breadcrumb = Handlebars.compile(source);

     $.post("/explore",{bucket:$('#ebucket').val()},function(data){
         log(data)

         // Render the bucket data
         var html = template({list: data.N});
         $('#expl').html(html)

        // Render the breadcrumbs
        var buckets = {}, parents = [];
        data.Buckets.forEach(function(bucket) {
            parents.push(bucket)
            buckets[bucket] = parents.join('/');
        })

        html = breadcrumb({buckets: buckets})
        $('#bdc').html(html)
     });
}

function loadBucketTable() {
    var source = $('#template').html();
    var template = Handlebars.compile(source);

    $.get("/buckets",{},function(data){
        var html    = template({list: data});
        $('#data').html(html)
    });
}

$(document).ready(function() {
    loadBucketTable();
    router.resolve();
});
