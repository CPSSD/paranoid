var selected = -1;

function loadSideBar() {
  $("#nav").empty();
  var items = [
  ];

  var starter = '<li><a href="#" ';
  if (selected == -1) {
    starter += 'class="selected" ';
  }
  starter += 'onclick="rowClicked(-1)">Create New</a></li>';
  items.push(starter);

  $.each(fileSystems, function(i, item) {
    var entry = '<li><a href="#" ';
    if (selected == i) {
      entry += 'class="selected" ';
    }
    if (item.mounted) {
      entry += 'id="mounted" ';
    }
    entry += 'onclick="rowClicked(' + i + ')">';
    entry += item.name + '</a></li>';
    items.push(entry);
  });

  $("#nav").append(items.join(' '));
}

function rowClicked(i) {
  if (i < -1 || i >= fileSystems.length) {
    i = -1;
  }
  selected = i;
  if (i == -1) {
    $(".content").load("html/form.html");
  } else {
    drawFileSystem(i);
  }
  loadSideBar();
}
