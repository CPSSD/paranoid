class MainActivity extends Activity{
  constructor(){
    super();
    this.layout = "main";
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
