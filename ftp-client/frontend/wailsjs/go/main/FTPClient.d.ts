// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {ftp} from '../models';

export function Connect(arg1:string,arg2:string,arg3:string):Promise<void>;

export function CreateFolder(arg1:string):Promise<void>;

export function Delete(arg1:string):Promise<void>;

export function Disconnect():Promise<void>;

export function Download(arg1:string,arg2:string):Promise<void>;

export function List(arg1:string):Promise<Array<ftp.Entry>>;

export function Upload(arg1:string,arg2:string):Promise<void>;
