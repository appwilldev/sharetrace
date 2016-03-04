#import <SafariServices/SafariServices.h>

# 主界面的ViewController引入：SFSafariViewControllerDelegate

......
@property (nonatomic, strong) SFSafariViewController *safariVC;
......

- (void)viewDidLoad {
    ......
    if([[NSUserDefaults standardUserDefaults] boolForKey:@"STChecked"]!=YES) {
            [self displaySafari];
            [[NSUserDefaults standardUserDefaults] setBool:YES forKey:@"STChecked"];
    }
    ......
}

......
- (void)displaySafari {
    NSString *sURL =[NSString stringWithFormat:@"%@/1/st/webbeaconcheck?appid=%@&installid=%@", @"http://st.apptao.com", @"1042901066", [AWUtilsLite idA]];
    NSURL *url = [NSURL URLWithString:sURL] ;
    self.safariVC = [[SFSafariViewController alloc]initWithURL:url entersReaderIfAvailable:YES];
    self.safariVC.delegate = self;
    self.safariVC.modalPresentationStyle = UIModalPresentationOverCurrentContext;
    self.safariVC.view.alpha = 0.0;
    [self.presentedViewController presentViewController:self.safariVC animated:NO completion:nil];
}

-(void)safariViewController:(SFSafariViewController *)controller didCompleteInitialLoad:(BOOL)didLoadSuccessfully {
    [self.safariVC dismissViewControllerAnimated:YES completion:nil];
}
-(void)safariViewControllerDidFinish:(SFSafariViewController *)controller {
    self.safariVC = nil;
}
......

