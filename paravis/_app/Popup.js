class Popup {
  constructor(title, layout){
    var self = this;
    console.log("Creating popup", title);

    self.m_title = title;
    self.m_layout = layout;
    self.m_element = document.createElement('popup');
    self.m_titleElement = document.createElement('title');
    self.m_contentElement = document.createElement('content');

    self.m_titleElement.textContent = this.m_title;
    $(self.m_contentElement).load('/res/layout/popup/'+self.m_layout+'.html');

    $(self.m_element).append(this.m_titleElement);
    $(self.m_element).append(this.m_contentElement);

    $(document).on('click', (e) => {
      console.log($(e.target).parents().is(self.m_element));
      if(!$(e.target).closest(self.m_element)){
        e.stopPropagation();
        self.hide();
      }
    });
  }

  hide(){
    $(this.m_element).remove();
  }

  show(){
    console.info('Showing popup')
    $(document.body).append(this.m_element);
  }
}
