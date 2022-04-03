import { Creatable } from "./creatable";

export interface Updatable extends Creatable {
  UpdatedAt: Date;
}
