class StartActivity extends Activity {
  constructor(){
    super();
    this.m_layout = "start";
  }

  onCreate(){
    loadLayout(this.m_layout);
  }
}

// Register the activity
activities[StartActivity.name] = StartActivity;
