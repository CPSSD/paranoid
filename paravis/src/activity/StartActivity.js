class StartActivity extends Activity {
  constructor(){
    super();
    this.m_layout = "start";
  }

  name(){
    return "StartActivity";
  }

  onCreate(){
    loadLayout(this.m_layout);
  }
}

// Register the activity
activities[StartActivity.name] = StartActivity;
