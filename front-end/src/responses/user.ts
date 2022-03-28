import { Updatable } from "./bases/updatable";
import { Circle } from "./circle";

export interface User extends Updatable {
  Name: string;
  Circles: Circle[];
}
