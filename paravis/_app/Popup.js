var console = require('console');

class Popup {
  constructor(title, layout, parent){
    var self = this;
    console.log("Creating popup", title);

    self.m_title = title;
    self.m_layout = layout;
    self.m_parent = parent;
    self.m_element = document.createElement('popup');
    self.m_titleElement = document.createElement('title');
    self.m_contentElement = document.createElement('content');

    self.m_titleElement.textContent = this.m_title;
    $(self.m_contentElement).load('/res/layout/popup/'+self.m_layout+'.html');

    $(self.m_element).append(this.m_titleElement);
    $(self.m_element).append(this.m_contentElement);

    $(document).on('click', (e) => {
      if(!$(self.m_element).is(e.target)){
        self.hide();
      }
    });
  }

  destructor(){
    delete(this.m_element);
  }

  show(){
    $(document.body).append(this.m_element);
  }
}

module.exports = Popup;
