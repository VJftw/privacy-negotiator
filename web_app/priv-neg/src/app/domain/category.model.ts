export class Category {
  name: string;
  personal: boolean;

  public static fromAPICategory(apiCat: APICategory): Category {
    return new Category(apiCat.name, apiCat.personal);
  }

  public static isIn(needle: Category, haystack: Category[]): boolean {
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

export class APICategory {
  name: string;
  personal: boolean;
}

export class CategorySelection {
  name: string;
  isActive: boolean;

  constructor(name: string, isActive = false) {
    this.name = name;
    this.isActive = isActive;
  }
}
