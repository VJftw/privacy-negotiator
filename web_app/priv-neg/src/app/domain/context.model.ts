export class Context {
  name: string;
  personal: boolean;

  public static fromAPIContext(apiCat: APIContext): Context {
    return new Context(apiCat.name, apiCat.personal);
  }

  public static isIn(needle: Context, haystack: Context[]): boolean {
    for (const hay of haystack) {
      if (needle.name === hay.name) {
        return true;
      }
    }
    return false;
  }

  constructor(c: string, p = false) {
    this.name = c;
    this.personal = p;
  }
}

export class APIContext {
  name: string;
  personal: boolean;
}

export class ContextSelection {
  name: string;
  isActive: boolean;

  constructor(name: string, isActive = false) {
    this.name = name;
    this.isActive = isActive;
  }
}
