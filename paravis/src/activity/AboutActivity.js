class AboutActivity extends Activity {
  constructor(){
    super();
    this.m_layout = "about";
  }

  onCreate(){
    loadLayout(this.m_layout);
  }
}

// Register the activity
activities[AboutActivity.name] = AboutActivity;
