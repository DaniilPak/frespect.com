export interface IRepository<T> {
  get(id: string): Promise<T>;
  getAll(): Promise<T[]>;
  create(item: T): Promise<T>;
  update(id: string, item: T): Promise<T>;
  delete(id: string): Promise<void>;
}
