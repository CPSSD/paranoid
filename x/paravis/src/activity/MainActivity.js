class MainActivity extends Activity{
  constructor(){
    super();
    this.layout = "main";
    this.db = new KeyValueDB(1, PermanentDB);
  }

  layout(){
    return this.layout;
  }

  onCreate() {
    alert(A.name);
  }
}

// All activies must be registered
activities[MainActivity.name] =  MainActivity;
