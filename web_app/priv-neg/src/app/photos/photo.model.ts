export class Photo {

  protected id;
  protected createdTime;
  protected from;
  protected images;

  public album;
  public privacy;

  protected sharedPosts;

  public static fromGraphAPI(obj): Photo {
    const p = new Photo();

    p.id = obj.id;
    p.createdTime = obj.created_time;
    p.from = obj.from;
    p.images = obj.images;
    p.album = obj.album;

    return p;
  }
}
