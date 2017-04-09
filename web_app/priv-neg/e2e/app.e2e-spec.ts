import { PrivNegPage } from './app.po';

describe('priv-neg App', () => {
  let page: PrivNegPage;

  beforeEach(() => {
    page = new PrivNegPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
