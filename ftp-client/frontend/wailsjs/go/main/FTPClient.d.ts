// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {context} from '../models';

export function Close():Promise<void>;

export function Dele(arg1:string,arg2:boolean):Promise<void>;

export function Dial(arg1:string):Promise<void>;

export function ListFiles(arg1:string):Promise<Array<string>>;

export function Login(arg1:string,arg2:string):Promise<void>;

export function MakeDir(arg1:string):Promise<void>;

export function REST_RETR(arg1:string,arg2:string,arg3:number,arg4:context.Context):Promise<void>;

export function RETR(arg1:string,arg2:string):Promise<void>;

export function STOR(arg1:string,arg2:string):Promise<void>;

export function SendCommand(arg1:string):Promise<string>;

export function SetAsciiMode():Promise<void>;

export function SetBinaryMode():Promise<void>;
