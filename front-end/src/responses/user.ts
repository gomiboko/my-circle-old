import { Circle } from "./circle";

export type User = {
  ID: number,
  Name: string,
  CreatedAt: Date,
  UpdatedAt: Date,
  Circles: Circle[],
}
