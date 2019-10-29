import { Selector } from 'testcafe';

export default {
   common:{
       mainMenu: Selector('.jss140.jss141'),
       shortlists: Selector('.jss251').find('div').withText('Shortlists')
   },
   shortlists:{
       search: Selector('input').withAttribute('placeholder', 'Search')
   }
};