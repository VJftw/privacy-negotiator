export class Category {
  name: string;

  constructor(c: string) {
    this.name = c;
  }
}

export class CategorySelection {
  name: string;
  isActive: boolean;

  constructor(name: string, isActive = false) {
    this.name = name;
    this.isActive = isActive;
  }
}
