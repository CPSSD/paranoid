const DIALOG_TIMEOUT = 2000 //ms

class Dialog {
  constructor(text, onClick, okButtonText, dismissButtonText){
    this.m_element = document.createElement('dialog');
    this.m_element.textContent = text;

    this.m_element.innerHtml += document.createElement('ok')

    setTimeout(() => {
      this.Remove();
    }, DIALOG_TIMEOUT);
  }

  Remove(){
    $(this.m_element).remove();
  }
}
